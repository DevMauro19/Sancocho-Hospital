package hospital

import "testing"

func TestOrderlyAndDoctorShareAttenderSlice(t *testing.T) {
	h := newTestHospital(t, 1)
	d := NewDoctor("D-1", "Karen", 40, "Neuro")
	o := NewOrderly("O-1", "Camillero", 30)
	if err := h.HireDoctor(d); err != nil {
		t.Fatal(err)
	}
	if err := h.HireAttender(o); err != nil {
		t.Fatal(err)
	}
	p1 := NewPatient("P-1", "Ana", 30, Severe)
	p2 := NewPatient("P-2", "Beto", 31, Mild)
	_ = h.AdmitPatient(p1)
	_ = h.AdmitPatient(p2)
	if _, err := h.HandleSleepAttack(p1, "cafeteria"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.HandleSleepAttack(p2, "radiology"); err == nil {
		t.Fatal("expected no-room error")
	}
	if len(d.MyEpisodes()) != 1 || len(o.MyEpisodes()) != 1 {
		t.Fatalf("doctor=%d orderly=%d, want 1/1", len(d.MyEpisodes()), len(o.MyEpisodes()))
	}
	// WakeUp libera el cuarto y el del pasillo lo recibe
	if err := p1.WakeUp(); err != nil {
		t.Fatal(err)
	}
	if p1.State() != Awake || p1.AssignedRoom() != nil {
		t.Fatal("p1 should be awake with no room")
	}
	if r, err := h.AssignRoom(p2); err != nil || r.Number() != 101 {
		t.Fatalf("p2 should get 101: %v %v", r, err)
	}
	if err := p1.WakeUp(); err == nil {
		t.Fatal("waking an awake patient should fail")
	}
	for _, rec := range h.History() {
		t.Log(rec.Summary())
	}
	t.Logf("State prints as: %v | level prints as: %v", p2.State(), p2.Level())
}