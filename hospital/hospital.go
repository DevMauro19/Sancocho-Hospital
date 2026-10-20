package hospital

import (
	"errors"
	"fmt"
	"strings"
)

// Errores centinela: permiten a quien llama distinguir el tipo de fallo con
// errors.Is, sin comparar textos.
var (
	// ErrNoRoomAvailable indica que ninguna habitación tiene cama libre.
	ErrNoRoomAvailable = errors.New("no room available")
	// ErrNoStaff indica que el hospital no tiene personal contratado.
	ErrNoStaff = errors.New("no staff hired")
)

// Hospital es la entidad central: coordina personal, pacientes, habitaciones
// e historial de episodios. Todos los campos son privados al paquete; desde
// fuera solo se llega a ellos con los métodos de este archivo.
type Hospital struct {
	name     string
	doctors  []*Doctor       // médicos contratados (campo mínimo del enunciado)
	staff    []Attender      // TODO el personal que puede atender: doctores, camilleros...
	patients []*Patient      // pacientes admitidos, en orden de admisión
	rooms    []*Room         // habitaciones registradas, en orden de registro
	history  []EpisodeRecord // todos los episodios de sueño registrados

	nextAttender int // índice por el que empieza el siguiente despacho (turno rotativo)
}

// NewHospital crea un hospital vacío. Falla si el nombre está en blanco.
func NewHospital(name string) (*Hospital, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("hospital name must not be empty")
	}
	return &Hospital{name: name}, nil
}

// Getters de solo lectura. Los que devuelven slices entregan una copia para
// que nadie de afuera modifique el estado interno.

// Name devuelve el nombre del hospital.
func (h *Hospital) Name() string { return h.name }

// Doctors devuelve una copia de la lista de médicos contratados.
func (h *Hospital) Doctors() []*Doctor { return append([]*Doctor(nil), h.doctors...) }

// Staff devuelve una copia de todo el personal (cualquier tipo de Attender).
func (h *Hospital) Staff() []Attender { return append([]Attender(nil), h.staff...) }

// Patients devuelve una copia de la lista de pacientes admitidos.
func (h *Hospital) Patients() []*Patient { return append([]*Patient(nil), h.patients...) }

// Rooms devuelve una copia de la lista de habitaciones.
func (h *Hospital) Rooms() []*Room { return append([]*Room(nil), h.rooms...) }

// History devuelve una copia del historial completo de episodios.
func (h *Hospital) History() []EpisodeRecord { return append([]EpisodeRecord(nil), h.history...) }

// HireDoctor contrata a un médico: queda en la lista de médicos y también en
// el personal general (staff), para que pueda ser despachado a un episodio.
func (h *Hospital) HireDoctor(d *Doctor) error {
	if d == nil {
		return errors.New("doctor is nil")
	}
	if err := h.hire(d); err != nil {
		return err
	}
	h.doctors = append(h.doctors, d)
	return nil
}

// HireAttender contrata a cualquier otro tipo de Attender (por ejemplo un
// Orderly). Los médicos se contratan con HireDoctor.
func (h *Hospital) HireAttender(a Attender) error {
	if a == nil {
		return errors.New("attender is nil")
	}
	return h.hire(a)
}

// hire agrega a staff validando que no haya otro con el mismo ID. Solo usa la
// interfaz Attender: no sabe si recibe un Doctor, un Orderly u otro tipo.
func (h *Hospital) hire(a Attender) error {
	for _, s := range h.staff {
		if s.ID() == a.ID() {
			return fmt.Errorf("staff member %s is already hired", a.ID())
		}
	}
	h.staff = append(h.staff, a)
	return nil
}

// AddRoom registra una habitación. Rechaza nil y números repetidos.
func (h *Hospital) AddRoom(r *Room) error {
	if r == nil {
		return errors.New("room is nil")
	}
	for _, existing := range h.rooms {
		if existing.Number() == r.Number() {
			return fmt.Errorf("room %d is already registered", r.Number())
		}
	}
	h.rooms = append(h.rooms, r)
	return nil
}

// AdmitPatient admite a un paciente. Rechaza nil e IDs repetidos.
func (h *Hospital) AdmitPatient(p *Patient) error {
	if p == nil {
		return errors.New("patient is nil")
	}
	for _, existing := range h.patients {
		if existing.ID() == p.ID() {
			return fmt.Errorf("patient %s is already admitted", p.ID())
		}
	}
	h.patients = append(h.patients, p)
	return nil
}

// hasPatient indica si el paciente fue admitido (compara punteros).
func (h *Hospital) hasPatient(p *Patient) bool {
	for _, existing := range h.patients {
		if existing == p {
			return true
		}
	}
	return false
}

// AssignRoom busca la primera habitación disponible, acuesta al paciente y
// actualiza su estado a AsleepInBed. Si no hay cama devuelve un error que
// envuelve ErrNoRoomAvailable y el paciente sigue en AsleepInHallway: no se
// inventa ninguna habitación ni se interrumpe el programa.
func (h *Hospital) AssignRoom(p *Patient) (*Room, error) {
	if p == nil {
		return nil, errors.New("patient is nil")
	}
	if !h.hasPatient(p) {
		return nil, fmt.Errorf("patient %s is not admitted", p.ID())
	}
	switch p.state {
	case Awake:
		return nil, fmt.Errorf("patient %s is awake: only sleeping patients need a bed", p.ID())
	case AsleepInBed:
		return nil, fmt.Errorf("patient %s is already in a bed", p.ID())
	}

	for _, r := range h.rooms {
		if !r.IsAvailable() {
			continue
		}
		if err := r.Occupy(p); err != nil {
			return nil, err
		}
		p.state = AsleepInBed
		p.assignedRoom = r
		return r, nil
	}
	return nil, fmt.Errorf("%w: patient %s stays in the hallway", ErrNoRoomAvailable, p.ID())
}

// RegisterEpisode guarda un episodio en el historial del hospital. Exige que
// el registro tenga paciente y que ese paciente esté admitido.
func (h *Hospital) RegisterEpisode(rec EpisodeRecord) error {
	if rec.patient == nil {
		return errors.New("episode has no patient")
	}
	if !h.hasPatient(rec.patient) {
		return fmt.Errorf("patient %s is not admitted", rec.patient.ID())
	}
	h.history = append(h.history, rec)
	return nil
}

// dispatch elige quién atiende el episodio. Recorre el personal desde el
// turno que sigue y prueba con cada Attender hasta que uno atiende sin error
// ("el que esté disponible"). Solo conoce la interfaz Attender: así el mismo
// código sirve para Doctor, Orderly o cualquier tipo futuro.
func (h *Hospital) dispatch(p *Patient, location string) (EpisodeRecord, error) {
	if len(h.staff) == 0 {
		return EpisodeRecord{}, ErrNoStaff
	}
	var failures []error
	for i := 0; i < len(h.staff); i++ {
		idx := (h.nextAttender + i) % len(h.staff)
		a := h.staff[idx]
		rec, err := a.Attend(p, location)
		if err == nil {
			h.nextAttender = (idx + 1) % len(h.staff) // el siguiente episodio empieza por el próximo
			return rec, nil
		}
		failures = append(failures, fmt.Errorf("%s (%s): %w", a.Name(), a.ID(), err))
	}
	return EpisodeRecord{}, errors.Join(failures...)
}

// HandleSleepAttack ejecuta el flujo completo de un ataque de sueño:
//  1. el paciente cae dormido en el pasillo (AsleepInHallway);
//  2. se le intenta asignar habitación;
//  3. se despacha a un Attender, que genera el EpisodeRecord;
//  4. el episodio se registra en el historial.
//
// El episodio se registra aunque no haya cama: el paciente tuvo el ataque y
// debe contar en el historial (lo necesita SevereReport). En ese caso el
// resultado es (nil, error que envuelve ErrNoRoomAvailable).
func (h *Hospital) HandleSleepAttack(p *Patient, location string) (*Room, error) {
	if p == nil {
		return nil, errors.New("patient is nil")
	}
	if !h.hasPatient(p) {
		return nil, fmt.Errorf("patient %s is not admitted", p.ID())
	}
	if p.state != Awake {
		return nil, fmt.Errorf("patient %s is already asleep", p.ID())
	}
	if strings.TrimSpace(location) == "" {
		return nil, errors.New("location must not be empty")
	}

	p.SufferSleepAttack(location)

	room, roomErr := h.AssignRoom(p)
	if roomErr != nil && !errors.Is(roomErr, ErrNoRoomAvailable) {
		return nil, roomErr
	}

	// El Attender se despacha después de asignar la cama para que el registro
	// del episodio incluya la habitación (o quede vacía si no hubo).
	rec, err := h.dispatch(p, location)
	if err != nil {
		return room, fmt.Errorf("dispatching an attender for %s: %w", p.ID(), err)
	}
	if err := h.RegisterEpisode(rec); err != nil {
		return room, err
	}
	return room, roomErr
}

// PatientsInHallway (consulta 5.1) devuelve los pacientes dormidos en un
// pasillo, en orden de admisión, para despachar a un camillero.
func (h *Hospital) PatientsInHallway() []*Patient {
	inHallway := make([]*Patient, 0)
	for _, p := range h.patients {
		if p.state == AsleepInHallway {
			inHallway = append(inHallway, p)
		}
	}
	return inHallway
}

// SevereReport (consulta 5.4) cruza dos fuentes: los pacientes con nivel Severe
// y el historial de episodios. Devuelve, por cada paciente grave, cuántos
// episodios tiene registrados (0 si no ha tenido ninguno).
func (h *Hospital) SevereReport() map[*Patient]int {
	report := make(map[*Patient]int)
	for _, p := range h.patients {
		if p.level == Severe {
			report[p] = 0
		}
	}
	for _, rec := range h.history {
		if _, isSevere := report[rec.patient]; isSevere {
			report[rec.patient]++
		}
	}
	return report
}
