package hospital

import (
	"errors"
	"fmt"
)

// RoomState es el estado de una habitación: disponible u ocupada.
type RoomState int

const (
	Available RoomState = iota // 0: hay al menos una cama libre
	Occupied                   // 1: el cuarto está lleno
)

// String permite imprimir el estado de forma legible (fmt lo usa automáticamente).
func (s RoomState) String() string {
	switch s {
	case Available:
		return "Available"
	case Occupied:
		return "Occupied"
	default:
		return "Unknown State"
	}
}

// Room representa una habitación del hospital. Los campos son privados:
// solo se leen o modifican a través de los métodos de este archivo.
type Room struct {
	number    int
	capacity  int
	state     RoomState
	occupants []*Patient // pacientes que están durmiendo en el cuarto
}

// NewRoom crea una habitación vacía. Falla si la capacidad no es positiva,
// porque un cuarto sin camas no tiene sentido en el hospital.
func NewRoom(number, capacity int) (*Room, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("room %d: capacity must be greater than 0, got %d", number, capacity)
	}
	return &Room{
		number:   number,
		capacity: capacity,
		state:    Available,
	}, nil
}

// Getters: dan acceso de solo lectura a los campos privados.
func (r *Room) Number() int      { return r.number }
func (r *Room) Capacity() int    { return r.capacity }
func (r *Room) State() RoomState { return r.state }

// Occupants devuelve una copia, para que nadie modifique el slice interno.
func (r *Room) Occupants() []*Patient {
	out := make([]*Patient, len(r.occupants))
	copy(out, r.occupants)
	return out
}

// IsAvailable indica si queda al menos una cama libre.
func (r *Room) IsAvailable() bool {
	return len(r.occupants) < r.capacity
}

// Occupy acuesta al paciente en el cuarto. Solo cambia el cuarto: el estado
// del paciente lo actualiza Hospital.AssignRoom.
func (r *Room) Occupy(p *Patient) error {
	if p == nil {
		return errors.New("patient is nil")
	}
	if r.contains(p) {
		return fmt.Errorf("patient %s is already in room %d", p.ID(), r.number)
	}
	if !r.IsAvailable() {
		return fmt.Errorf("room %d is full", r.number)
	}
	r.occupants = append(r.occupants, p)
	r.updateState()
	return nil
}

// Release saca al paciente del cuarto y deja la cama libre.
func (r *Room) Release(p *Patient) error {
	if p == nil {
		return errors.New("patient is nil")
	}
	for i, occ := range r.occupants {
		if occ == p {
			// Quita el elemento i sin dejar huecos en el slice.
			r.occupants = append(r.occupants[:i], r.occupants[i+1:]...)
			r.updateState()
			return nil
		}
	}
	return fmt.Errorf("patient %s is not in room %d", p.ID(), r.number)
}

// contains revisa si el paciente ya está en el cuarto (compara punteros).
func (r *Room) contains(p *Patient) bool {
	for _, occ := range r.occupants {
		if occ == p {
			return true
		}
	}
	return false
}

// updateState recalcula el estado según los ocupantes. Es el único lugar
// donde cambia r.state, así nunca queda "Occupied" con camas libres.
func (r *Room) updateState() {
	if r.IsAvailable() {
		r.state = Available
	} else {
		r.state = Occupied
	}
}
