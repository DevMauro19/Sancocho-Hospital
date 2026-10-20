package hospital

import "errors"

type PatientState int

const (
	Awake PatientState = iota
	AsleepInHallway
	AsleepInBed
)

func (s PatientState) String() string {
	switch s {
	case Awake:
		return "Awake"
	case AsleepInHallway:
		return "Asleep in Hallway"
	case AsleepInBed:
		return "Asleep in Bed"
	default:
		return "Unknown State"
	}
}

type NarcolepsyLevel int

const (
	Mild NarcolepsyLevel = iota
	Moderate
	Severe
)

func (l NarcolepsyLevel) String() string {
	switch l {
	case Mild:
		return "Mild"
	case Moderate:
		return "Moderate"
	case Severe:
		return "Severe"
	default:
		return "Unknown Level"
	}
}

type Person struct {
	id   string
	name string
	age  int
}

func (p Person) ID() string   { return p.id }
func (p Person) Name() string { return p.name }
func (p Person) Age() int     { return p.age }

type Patient struct {
	Person
	level           NarcolepsyLevel
	state           PatientState
	currentLocation string
	assignedDoctor  *Doctor
	assignedRoom    *Room
}

func NewPatient(id, name string, age int, level NarcolepsyLevel) *Patient {
	return &Patient{Person: Person{id: id, name: name, age: age}, level: level}
}

func (p *Patient) State() PatientState    { return p.state }
func (p *Patient) Level() NarcolepsyLevel { return p.level }
func (p *Patient) Location() string       { return p.currentLocation }
func (p *Patient) AssignedRoom() *Room    { return p.assignedRoom }

// SufferSleepAttack: el paciente cae dormido en pasillo.
func (p *Patient) SufferSleepAttack(location string) {
	p.currentLocation = location
	p.state = AsleepInHallway
}

// WakeUp despierta al paciente y libera su cuarto si tenía uno.
func (p *Patient) WakeUp() error {
	if p.state == Awake {
		return errors.New("patient is already awake")
	}
	if p.assignedRoom != nil {
		if err := p.assignedRoom.Release(p); err != nil {
			return err
		}
		p.assignedRoom = nil
	}
	p.state = Awake
	return nil
}

type Attender interface {
	ID() string
	Name() string
	Attend(p *Patient, location string) (EpisodeRecord, error)
}