package main

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"costenitos/hospital"
)

func main() {
	// --- 1. Hospital, personal, habitaciones y pacientes -------------------
	h, err := hospital.NewHospital("Hospital de los Costeños con Narcolepsia")
	check(err)

	doctors := []*hospital.Doctor{
		hospital.NewDoctor("D-001", "Karen Ospina", 41, "Neurología"),
		hospital.NewDoctor("D-002", "Andrés Caballero", 47, "Medicina del sueño"),
	}
	for _, d := range doctors {
		check(h.HireDoctor(d))
	}
	orderly := hospital.NewOrderly("O-001", "Jairo Camacho", 33)
	check(h.HireAttender(orderly))

	for number := 101; number <= 103; number++ {
		room, err := hospital.NewRoom(number, 1)
		check(err)
		check(h.AddRoom(room))
	}

	patients := []*hospital.Patient{
		hospital.NewPatient("P-001", "Ana Torres", 34, hospital.Mild),
		hospital.NewPatient("P-002", "Yeimy Padilla", 29, hospital.Severe),
		hospital.NewPatient("P-003", "Carlos Mendoza", 52, hospital.Moderate),
		hospital.NewPatient("P-004", "Wilfrido Berrio", 45, hospital.Severe),
		hospital.NewPatient("P-005", "Luz Marina Pérez", 61, hospital.Mild),
	}
	for _, p := range patients {
		check(h.AdmitPatient(p))
	}

	fmt.Printf("== %s ==\n", h.Name())
	fmt.Println("Staff (a []Attender holding different concrete types):")
	for _, a := range h.Staff() {
		fmt.Printf("  %-6s %-18s %T\n", a.ID(), a.Name(), a)
	}
	fmt.Printf("Rooms: %d | Patients admitted: %d\n\n", len(h.Rooms()), len(h.Patients()))

	// --- 2. Cuatro ataques de sueño: tres reciben cuarto, el cuarto falla ---
	attacks := []struct {
		patient  *hospital.Patient
		location string
	}{
		{patients[0], "cafeteria"},
		{patients[1], "radiology queue"},
		{patients[2], "pharmacy"},
		{patients[3], "hallway 2, second floor"},
	}

	fmt.Println("== Sleep attacks ==")
	var noRoomErr error
	for _, a := range attacks {
		room, err := h.HandleSleepAttack(a.patient, a.location)
		switch {
		case err == nil:
			fmt.Printf("  %s %-16s fell asleep at %-24s -> room %d\n",
				a.patient.ID(), a.patient.Name(), a.location, room.Number())
		case errors.Is(err, hospital.ErrNoRoomAvailable):
			noRoomErr = err
			fmt.Printf("  %s %-16s fell asleep at %-24s -> %v\n",
				a.patient.ID(), a.patient.Name(), a.location, err)
		default:
			check(err)
		}
	}
	// Consulta 5.1 tomada justo después de los ataques (antes de que nadie despierte).
	hallwayAfterAttacks := h.PatientsInHallway()

	// --- 3. Un paciente despierta, libera su cuarto y el del pasillo lo recibe
	fmt.Println("\n== A patient wakes up ==")
	check(patients[0].WakeUp())
	fmt.Printf("  %s %s woke up and released the room\n", patients[0].ID(), patients[0].Name())

	reassigned, err := h.AssignRoom(patients[3])
	check(err)
	fmt.Printf("  %s %s is now in room %d\n\n", patients[3].ID(), patients[3].Name(), reassigned.Number())

	// --- 4. Las cuatro consultas de la sección 5 ----------------------------
	fmt.Println("== 5.1 Patients asleep in hallway (after the four attacks) ==")
	printPatients(hallwayAfterAttacks)
	fmt.Println("== 5.1 Patients asleep in hallway (now) ==")
	printPatients(h.PatientsInHallway())

	for _, d := range doctors {
		fmt.Printf("== 5.2 Episodes attended by Dr. %s ==\n", d.Name())
		printEpisodes(d.MyEpisodes())
	}
	fmt.Printf("== 5.2 Episodes attended by orderly %s ==\n", orderly.Name())
	printEpisodes(orderly.MyEpisodes())

	fmt.Printf("== 5.3 Assigning room to %s ==\n", attacks[3].patient.ID())
	fmt.Printf("  first attempt: %v\n", noRoomErr)
	fmt.Printf("  after %s woke up: assigned to room %d\n\n", patients[0].ID(), reassigned.Number())

	fmt.Println("== 5.4 Severity report (Severe patients) ==")
	report := h.SevereReport()
	severe := make([]*hospital.Patient, 0, len(report))
	for p := range report {
		severe = append(severe, p)
	}
	// Los mapas no tienen orden en Go: se ordena por ID para una salida estable.
	sort.Slice(severe, func(i, j int) bool { return severe[i].ID() < severe[j].ID() })
	for _, p := range severe {
		fmt.Printf("  %s %-16s %d episodes today\n", p.ID(), p.Name(), report[p])
	}
}

// printPatients imprime una lista de pacientes con su ubicación.
func printPatients(ps []*hospital.Patient) {
	if len(ps) == 0 {
		fmt.Println("  (none)")
	}
	for _, p := range ps {
		fmt.Printf("  %s %-16s (%s)\n", p.ID(), p.Name(), p.Location())
	}
	fmt.Println()
}

// printEpisodes imprime los episodios usando Summary().
func printEpisodes(records []hospital.EpisodeRecord) {
	if len(records) == 0 {
		fmt.Println("  (none)")
	}
	for _, rec := range records {
		fmt.Printf("  %s\n", rec.Summary())
	}
	fmt.Println()
}

// check termina el programa si falla algo del montaje del escenario.
func check(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
