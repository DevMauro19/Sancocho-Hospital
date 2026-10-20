package hospital

import (
	"errors"
	"testing"
)

// --- Ayudas de prueba -------------------------------------------------------

// fakeAttender es un Attender de mentira para probar el despacho sin depender
// de Doctor ni de Orderly. Cuenta cuántos episodios atiende y puede forzarse a
// fallar con failWith.
type fakeAttender struct {
	id       string
	name     string
	attended int
	failWith error
}

func (f *fakeAttender) ID() string   { return f.id }
func (f *fakeAttender) Name() string { return f.name }

func (f *fakeAttender) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if f.failWith != nil {
		return EpisodeRecord{}, f.failWith
	}
	f.attended++
	return EpisodeRecord{patient: p, location: location}, nil
}

// fakeOrderly es un SEGUNDO tipo concreto: embebe fakeAttender (composición) y
// por eso también cumple Attender. Permite probar el []Attender con dos tipos.
type fakeOrderly struct {
	fakeAttender
}

func newTestPatient(id string, level NarcolepsyLevel) *Patient {
	return &Patient{Person: Person{id: id, name: "Paciente " + id, age: 40}, level: level}
}

// newTestHospital crea un hospital con n habitaciones de capacidad 1 y el
// personal indicado.
func newTestHospital(t *testing.T, rooms int, staff ...Attender) *Hospital {
	t.Helper()
	h, err := NewHospital("Hospital de prueba")
	if err != nil {
		t.Fatalf("NewHospital: %v", err)
	}
	for i := 1; i <= rooms; i++ {
		r, err := NewRoom(100+i, 1)
		if err != nil {
			t.Fatalf("NewRoom: %v", err)
		}
		if err := h.AddRoom(r); err != nil {
			t.Fatalf("AddRoom: %v", err)
		}
	}
	for _, a := range staff {
		if err := h.HireAttender(a); err != nil {
			t.Fatalf("HireAttender: %v", err)
		}
	}
	return h
}

// admitAsleepInHallway admite al paciente y lo deja dormido en el pasillo sin
// pasar por HandleSleepAttack, para probar AssignRoom de forma aislada.
func admitAsleepInHallway(t *testing.T, h *Hospital, p *Patient) {
	t.Helper()
	if err := h.AdmitPatient(p); err != nil {
		t.Fatalf("AdmitPatient: %v", err)
	}
	p.state = AsleepInHallway
}

// --- Tests obligatorios del enunciado ----------------------------------------

// 1) Asignación de habitación cuando SÍ hay una libre.
func TestAssignRoomWithFreeRoom(t *testing.T) {
	h := newTestHospital(t, 1)
	p := newTestPatient("P-001", Mild)
	admitAsleepInHallway(t, h, p)

	room, err := h.AssignRoom(p)
	if err != nil {
		t.Fatalf("AssignRoom returned an unexpected error: %v", err)
	}
	if room == nil || room.Number() != 101 {
		t.Fatalf("expected room 101, got %v", room)
	}
	if p.state != AsleepInBed {
		t.Errorf("patient state = %v, want AsleepInBed", p.state)
	}
	if p.assignedRoom != room {
		t.Error("patient should keep a reference to its room")
	}
	if room.State() != Occupied {
		t.Errorf("room state = %v, want Occupied", room.State())
	}
}

// 2) Fallo cuando NO hay habitación libre: error, sin pánico, el paciente
// sigue en el pasillo y no se inventa un cuarto.
func TestAssignRoomWithoutFreeRoom(t *testing.T) {
	h := newTestHospital(t, 1)
	first := newTestPatient("P-001", Mild)
	second := newTestPatient("P-002", Mild)
	admitAsleepInHallway(t, h, first)
	admitAsleepInHallway(t, h, second)

	if _, err := h.AssignRoom(first); err != nil {
		t.Fatalf("first patient should get the only room: %v", err)
	}

	room, err := h.AssignRoom(second)
	if !errors.Is(err, ErrNoRoomAvailable) {
		t.Fatalf("error = %v, want ErrNoRoomAvailable", err)
	}
	if room != nil {
		t.Errorf("no room should be returned, got %v", room)
	}
	if second.state != AsleepInHallway {
		t.Errorf("patient state = %v, want AsleepInHallway", second.state)
	}
	if second.assignedRoom != nil {
		t.Error("patient must not have a room assigned")
	}
}

// 3) Una de las consultas: PatientsInHallway devuelve solo los del pasillo.
func TestPatientsInHallway(t *testing.T) {
	h := newTestHospital(t, 0)
	hallway := newTestPatient("P-001", Mild)
	inBed := newTestPatient("P-002", Moderate)
	awake := newTestPatient("P-003", Severe)
	for _, p := range []*Patient{hallway, inBed, awake} {
		if err := h.AdmitPatient(p); err != nil {
			t.Fatalf("AdmitPatient: %v", err)
		}
	}
	hallway.state = AsleepInHallway
	inBed.state = AsleepInBed

	got := h.PatientsInHallway()
	if len(got) != 1 || got[0] != hallway {
		t.Fatalf("PatientsInHallway() = %v, want only P-001", got)
	}
}

// --- Tests adicionales --------------------------------------------------------

// SevereReport: solo pacientes graves, con su número de episodios (0 incluido).
func TestSevereReport(t *testing.T) {
	h := newTestHospital(t, 0)
	severeBusy := newTestPatient("P-001", Severe)
	severeQuiet := newTestPatient("P-002", Severe)
	mild := newTestPatient("P-003", Mild)
	for _, p := range []*Patient{severeBusy, severeQuiet, mild} {
		if err := h.AdmitPatient(p); err != nil {
			t.Fatalf("AdmitPatient: %v", err)
		}
	}
	for _, p := range []*Patient{severeBusy, severeBusy, mild} {
		if err := h.RegisterEpisode(EpisodeRecord{patient: p, location: "cafeteria"}); err != nil {
			t.Fatalf("RegisterEpisode: %v", err)
		}
	}

	report := h.SevereReport()
	if len(report) != 2 {
		t.Fatalf("report has %d entries, want 2 (only severe patients)", len(report))
	}
	if report[severeBusy] != 2 {
		t.Errorf("severeBusy episodes = %d, want 2", report[severeBusy])
	}
	if count, ok := report[severeQuiet]; !ok || count != 0 {
		t.Errorf("severeQuiet = (%d, %v), want (0, true)", count, ok)
	}
	if _, ok := report[mild]; ok {
		t.Error("a non-severe patient must not appear in the report")
	}
}

// El despacho trabaja con []Attender de DOS tipos concretos distintos y reparte
// por turnos. Aunque el segundo paciente no consigue cama, su episodio queda
// registrado.
func TestHandleSleepAttackDispatchesPolymorphically(t *testing.T) {
	doctor := &fakeAttender{id: "D-001", name: "Dra. Ospina"}
	orderly := &fakeOrderly{fakeAttender{id: "O-001", name: "Camillero Pérez"}}
	h := newTestHospital(t, 1, doctor, orderly)

	p1 := newTestPatient("P-001", Mild)
	p2 := newTestPatient("P-002", Severe)
	for _, p := range []*Patient{p1, p2} {
		if err := h.AdmitPatient(p); err != nil {
			t.Fatalf("AdmitPatient: %v", err)
		}
	}

	room, err := h.HandleSleepAttack(p1, "cafeteria")
	if err != nil || room == nil {
		t.Fatalf("first attack: room=%v err=%v, want a room and no error", room, err)
	}

	room, err = h.HandleSleepAttack(p2, "radiology queue")
	if !errors.Is(err, ErrNoRoomAvailable) || room != nil {
		t.Fatalf("second attack: room=%v err=%v, want nil room and ErrNoRoomAvailable", room, err)
	}
	if p2.state != AsleepInHallway {
		t.Errorf("p2 state = %v, want AsleepInHallway", p2.state)
	}

	if len(h.History()) != 2 {
		t.Errorf("history has %d episodes, want 2 (also the one without a bed)", len(h.History()))
	}
	if doctor.attended != 1 || orderly.attended != 1 {
		t.Errorf("attended: doctor=%d orderly=%d, want 1 and 1 (turnos)", doctor.attended, orderly.attended)
	}
}

// Si quien toca el turno falla, el despacho pasa al siguiente; si todos fallan
// devuelve error y no registra nada.
func TestDispatchFailover(t *testing.T) {
	busy := &fakeAttender{id: "D-001", name: "Ocupado", failWith: errors.New("busy")}
	free := &fakeAttender{id: "D-002", name: "Libre"}
	h := newTestHospital(t, 1, busy, free)
	p := newTestPatient("P-001", Mild)
	if err := h.AdmitPatient(p); err != nil {
		t.Fatalf("AdmitPatient: %v", err)
	}

	if _, err := h.HandleSleepAttack(p, "hall"); err != nil {
		t.Fatalf("HandleSleepAttack: %v", err)
	}
	if free.attended != 1 || busy.attended != 0 {
		t.Errorf("attended: free=%d busy=%d, want 1 and 0", free.attended, busy.attended)
	}

	// Ahora fallan todos.
	free.failWith = errors.New("busy too")
	if err := p.assignedRoom.Release(p); err != nil {
		t.Fatalf("Release: %v", err)
	}
	p.assignedRoom = nil
	p.state = Awake
	if _, err := h.HandleSleepAttack(p, "hall 2"); err == nil {
		t.Fatal("expected an error when every attender fails")
	}
	if len(h.History()) != 1 {
		t.Errorf("history has %d episodes, want 1 (the failed one is not registered)", len(h.History()))
	}
}

// Sin personal contratado no se puede despachar a nadie.
func TestHandleSleepAttackWithoutStaff(t *testing.T) {
	h := newTestHospital(t, 1)
	p := newTestPatient("P-001", Mild)
	if err := h.AdmitPatient(p); err != nil {
		t.Fatalf("AdmitPatient: %v", err)
	}
	if _, err := h.HandleSleepAttack(p, "hall"); !errors.Is(err, ErrNoStaff) {
		t.Fatalf("error = %v, want ErrNoStaff", err)
	}
}

// Validaciones de entrada: nil, duplicados y estados inválidos.
func TestInputValidation(t *testing.T) {
	h := newTestHospital(t, 0)
	p := newTestPatient("P-001", Mild)

	if err := h.AdmitPatient(nil); err == nil {
		t.Error("AdmitPatient(nil) should fail")
	}
	if err := h.AdmitPatient(p); err != nil {
		t.Fatalf("AdmitPatient: %v", err)
	}
	if err := h.AdmitPatient(newTestPatient("P-001", Severe)); err == nil {
		t.Error("admitting a duplicate ID should fail")
	}
	if err := h.HireDoctor(nil); err == nil {
		t.Error("HireDoctor(nil) should fail")
	}
	if err := h.HireAttender(&fakeAttender{id: "A-1"}); err != nil {
		t.Fatalf("HireAttender: %v", err)
	}
	if err := h.HireAttender(&fakeAttender{id: "A-1"}); err == nil {
		t.Error("hiring a duplicate ID should fail")
	}
	if _, err := h.AssignRoom(p); err == nil {
		t.Error("AssignRoom for an awake patient should fail")
	}
	if err := h.RegisterEpisode(EpisodeRecord{}); err == nil {
		t.Error("RegisterEpisode without patient should fail")
	}
	if _, err := NewHospital("   "); err == nil {
		t.Error("NewHospital with a blank name should fail")
	}
}
