package hospital

import (
	"errors"
	"time"
)

// Orderly (camillero): mueve a un paciente dormido, pero NO diagnostica
// (no tiene DiagnosePatient). Cumple Attender igual que Doctor.
type Orderly struct {
	Person
	episodes []EpisodeRecord
}

var _ Attender = (*Orderly)(nil)

func NewOrderly(id, name string, age int) *Orderly {
	return &Orderly{Person: Person{id: id, name: name, age: age}}
}

func (o *Orderly) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if p == nil {
		return EpisodeRecord{}, errors.New("orderly: patient is nil")
	}
	if p.state == Awake {
		return EpisodeRecord{}, errors.New("orderly: patient is awake, nothing to move")
	}
	rec := EpisodeRecord{
		id:       newEpisodeID(),
		dateTime: time.Now(),
		patient:  p,
		attender: o,
		location: location,
		room:     p.assignedRoom,
	}
	o.episodes = append(o.episodes, rec)
	return rec, nil
}

func (o *Orderly) MyEpisodes() []EpisodeRecord {
	out := make([]EpisodeRecord, len(o.episodes))
	copy(out, o.episodes)
	return out
}