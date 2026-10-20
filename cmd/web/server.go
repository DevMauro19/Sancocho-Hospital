package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"

	"costenitos/hospital"
)

type app struct {
	mu       sync.Mutex
	h        *hospital.Hospital
	patients map[string]*hospital.Patient
}

type patientView struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Severity string `json:"severity"`
	State    string `json:"state"`
	Location string `json:"location"`
	Room     *int   `json:"room"`
}
type roomView struct {
	Number    int      `json:"number"`
	State     string   `json:"state"`
	Capacity  int      `json:"capacity"`
	Occupants []string `json:"occupants"`
}
type episodeView struct {
	Summary string `json:"summary"`
}
type stateView struct {
	Hospital string        `json:"hospital"`
	Patients []patientView `json:"patients"`
	Rooms    []roomView    `json:"rooms"`
	Hallway  []string      `json:"hallway"`
	Severe   []patientView `json:"severe"`
	Episodes []episodeView `json:"episodes"`
}

func newApp() (*app, error) {
	a := &app{}
	if err := a.reset(); err != nil {
		return nil, err
	}
	return a, nil
}
func (a *app) reset() error {
	h, err := hospital.NewHospital("Hospital de los Costeños con Narcolepsia")
	if err != nil {
		return err
	}
	for _, d := range []*hospital.Doctor{hospital.NewDoctor("D-001", "Karen Ospina", 41, "Neurología"), hospital.NewDoctor("D-002", "Andrés Caballero", 47, "Medicina del sueño")} {
		if err := h.HireDoctor(d); err != nil {
			return err
		}
	}
	if err := h.HireAttender(hospital.NewOrderly("O-001", "Jairo Camacho", 33)); err != nil {
		return err
	}
	for n := 101; n <= 103; n++ {
		r, e := hospital.NewRoom(n, 1)
		if e != nil {
			return e
		}
		if e = h.AddRoom(r); e != nil {
			return e
		}
	}
	ps := []*hospital.Patient{
		hospital.NewPatient("P-001", "Ana Torres", 34, hospital.Mild),
		hospital.NewPatient("P-002", "Yeimy Padilla", 29, hospital.Severe),
		hospital.NewPatient("P-003", "Carlos Mendoza", 52, hospital.Moderate),
		hospital.NewPatient("P-004", "Wilfrido Berrio", 45, hospital.Severe),
		hospital.NewPatient("P-005", "Luz Marina Pérez", 61, hospital.Mild),
	}
	pm := map[string]*hospital.Patient{}
	for _, p := range ps {
		if err := h.AdmitPatient(p); err != nil {
			return err
		}
		pm[p.ID()] = p
	}
	a.h = h
	a.patients = pm
	return nil
}
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	writeJSON(w, map[string]string{"error": err.Error()})
}
func (a *app) stateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", 405)
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	writeJSON(w, a.snapshot())
}
func (a *app) snapshot() stateView {
	v := stateView{
		Hospital: a.h.Name(),
		Patients: []patientView{},
		Rooms:    []roomView{},
		Hallway:  []string{},
		Severe:   []patientView{},
		Episodes: []episodeView{},
	}
	for _, p := range a.h.Patients() {
		pv := patientView{ID: p.ID(), Name: p.Name(), Age: p.Age(), Severity: p.Level().String(), State: p.State().String(), Location: p.Location()}
		if room := p.AssignedRoom(); room != nil {
			n := room.Number()
			pv.Room = &n
		}
		v.Patients = append(v.Patients, pv)
		if p.Level() == hospital.Severe {
			v.Severe = append(v.Severe, pv)
		}
	}
	for _, r := range a.h.Rooms() {
		rv := roomView{Number: r.Number(), State: r.State().String(), Capacity: r.Capacity(), Occupants: []string{}}
		for _, p := range r.Occupants() {
			rv.Occupants = append(rv.Occupants, p.Name())
		}
		v.Rooms = append(v.Rooms, rv)
	}
	for _, p := range a.h.PatientsInHallway() {
		v.Hallway = append(v.Hallway, p.Name())
	}
	for _, d := range a.h.Doctors() {
		for _, e := range d.MyEpisodes() {
			v.Episodes = append(v.Episodes, episodeView{e.Summary()})
		}
	}
	for _, s := range a.h.Staff() {
		if o, ok := s.(*hospital.Orderly); ok {
			for _, e := range o.MyEpisodes() {
				v.Episodes = append(v.Episodes, episodeView{e.Summary()})
			}
		}
	}
	sort.Slice(v.Episodes, func(i, j int) bool { return v.Episodes[i].Summary < v.Episodes[j].Summary })
	return v
}

type actionRequest struct {
	PatientID string `json:"patientId"`
	Location  string `json:"location"`
}

type customPatientRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Severity string `json:"severity"`
}

func parseSeverity(level string) (hospital.NarcolepsyLevel, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "", "mild":
		return hospital.Mild, nil
	case "moderate":
		return hospital.Moderate, nil
	case "severe":
		return hospital.Severe, nil
	default:
		return hospital.Mild, fmt.Errorf("severity %q is not valid: use Mild, Moderate or Severe", level)
	}
}

func (a *app) nextPatientID() string {
	used := map[string]bool{}
	for _, p := range a.h.Patients() {
		used[p.ID()] = true
	}
	for i := 1; ; i++ {
		id := fmt.Sprintf("P-%03d", i)
		if !used[id] {
			return id
		}
	}
}

func decodeAction(w http.ResponseWriter, r *http.Request) (actionRequest, bool) {
	var req actionRequest
	if r.Method != "POST" {
		http.Error(w, "method not allowed", 405)
		return req, false
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, err)
		return req, false
	}
	return req, true
}
func (a *app) addPatientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req customPatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, errors.New("name is required"))
		return
	}
	if req.Age < 0 {
		writeError(w, http.StatusBadRequest, errors.New("age must be zero or greater"))
		return
	}
	level, err := parseSeverity(req.Severity)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	id := strings.TrimSpace(req.ID)
	if id == "" {
		a.mu.Lock()
		id = a.nextPatientID()
		a.mu.Unlock()
	}
	if strings.TrimSpace(id) == "" {
		writeError(w, http.StatusBadRequest, errors.New("patient id cannot be empty"))
		return
	}
	p := hospital.NewPatient(id, req.Name, req.Age, level)
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.patients[id] != nil {
		writeError(w, http.StatusConflict, fmt.Errorf("patient %s already exists", id))
		return
	}
	if err := a.h.AdmitPatient(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	a.patients[id] = p
	writeJSON(w, map[string]any{"state": a.snapshot(), "message": "Paciente agregado: " + p.Name()})
}

func (a *app) deletePatientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct { PatientID string `json:"patientId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.PatientID == "" {
		writeError(w, http.StatusBadRequest, errors.New("patientId is required"))
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	p := a.patients[req.PatientID]
	if p == nil {
		writeError(w, http.StatusNotFound, errors.New("paciente no encontrado"))
		return
	}
	if err := a.h.DischargePatient(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	delete(a.patients, req.PatientID)
	writeJSON(w, map[string]any{"state": a.snapshot(), "message": "Paciente dado de alta: " + p.Name()})
}

func (a *app) attackHandler(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeAction(w, r)
	if !ok {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	p := a.patients[req.PatientID]
	if p == nil {
		writeError(w, 404, errors.New("paciente no encontrado"))
		return
	}
	_, err := a.h.HandleSleepAttack(p, req.Location)
	if err != nil && !errors.Is(err, hospital.ErrNoRoomAvailable) {
		writeError(w, 400, err)
		return
	}
	writeJSON(w, map[string]any{"state": a.snapshot(), "message": func() string {
		if err != nil {
			return "Ataque registrado; no hay cama disponible y el paciente permanece en el pasillo."
		}
		return "Ataque registrado y habitación asignada."
	}()})
}
func (a *app) wakeHandler(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeAction(w, r)
	if !ok {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	p := a.patients[req.PatientID]
	if p == nil {
		writeError(w, 404, errors.New("paciente no encontrado"))
		return
	}
	if err := p.WakeUp(); err != nil {
		writeError(w, 400, err)
		return
	}
	writeJSON(w, map[string]any{"state": a.snapshot(), "message": "Paciente despierto; la cama quedó libre."})
}
func (a *app) assignHandler(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeAction(w, r)
	if !ok {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	p := a.patients[req.PatientID]
	if p == nil {
		writeError(w, 404, errors.New("paciente no encontrado"))
		return
	}
	room, err := a.h.AssignRoom(p)
	if err != nil {
		writeError(w, 400, err)
		return
	}
	writeJSON(w, map[string]any{"state": a.snapshot(), "message": "Habitación asignada: " + roomNumber(room.Number())})
}

func (a *app) patientDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/patients/")
	if id == "" || id == "/api/patients" {
		writeError(w, http.StatusBadRequest, errors.New("patientId is required"))
		return
	}
	payload := struct { PatientID string `json:"patientId"` }{}
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil && !errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, err)
			return
		}
	}
	if payload.PatientID != "" {
		id = payload.PatientID
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	p := a.patients[id]
	if p == nil {
		writeError(w, http.StatusNotFound, errors.New("paciente no encontrado"))
		return
	}
	if err := a.h.DischargePatient(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	delete(a.patients, id)
	writeJSON(w, map[string]any{"state": a.snapshot(), "message": "Paciente dado de alta: " + p.Name()})
}

func roomNumber(n int) string { return fmt.Sprintf("%d", n) }
func (a *app) resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", 405)
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.reset(); err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, map[string]any{"state": a.snapshot(), "message": "Escenario reiniciado."})
}
