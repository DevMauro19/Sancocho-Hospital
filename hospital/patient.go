package hospital

type PatientState int //Tipo de dato personalizado basado en enteros (int)
type narcolepsyLevel int //Tipo de dato personalizado basado en enteros (int)

const (
    Awake= iota //mecanismo de enumeración automática, asigna valores enteros consecutivos a las constantes
    AsleepInHallway // 0 (Awake), 1 (AsleepInHallway), 2 (AsleepInBed)
    AsleepInBed
)
const(
    Mild =iota
    Moderate
    Severe
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

//Funciona de manera similar a una clase en otros lenguajes de programación, encapsulando datos y comportamientos relacionados en una sola entidad.
type Person struct {
    id string
    name string
    age int
}

//Funcionan de manera similar a los getters en otros lenguajes
func (p Person) ID() string {return p.id}
func (p Person) Name() string {return p.name}

//Representacion de un paciente
type Patient struct{
    Person // Embedding the Person struct
    level NarcolepsyLevel
    state PatientState
    currentLocation string //Pendiente
    assignedDoctor *Doctor //Pendiente
}

type Attender interface {
    ID() string
    Name () string
    Attend(p *Patient, location string) (EpisodeRecord, error)
}