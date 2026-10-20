package hospital

import (
	"errors"
	"time"
)

type Doctor struct {
	Person                    // Embedding the Person struct
	specialty string          // especialidad del doctor
	patients  []*Patient      // Pacientes bajo su cuidado
	episodes  []EpisodeRecord // Episodios de sueño en los que intervino este doctor
}

var _ Attender = (*Doctor)(nil) // Verifica que Doctor implemente la interfaz Attender

func NewDoctor(id, name string, age int, specialty string) *Doctor {
	return &Doctor{
		Person: Person{
			id:   id,
			name: name,
			age:  age,
		},
		specialty: specialty,
	}
}

// Specialty devuelve la especialidad del doctor.
func (d *Doctor) Specialty() string { return d.specialty }

// DiagnosePatient asigna al paciente a este doctor.
func (d *Doctor) DiagnosePatient(p *Patient) error {
	if p == nil {
		return errors.New("patient is nil")
	}
	if p.assignedDoctor != nil && p.assignedDoctor != d {
		return errors.New("patient is already assigned to another doctor")
	}
	if p.assignedDoctor == nil {
		p.assignedDoctor = d
		d.patients = append(d.patients, p)
	}
	return nil
}

// Attend implementa Attender: atiende un episodio y lo registra en el
// historial del doctor. El historial global lo guarda Hospital, y el paciente
// no guarda episodios (para no duplicar el mismo dato en tres sitios).
func (d *Doctor) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if p == nil {
		return EpisodeRecord{}, errors.New("patient is nil")
	}
	if p.state == Awake {
		return EpisodeRecord{}, errors.New("patient is awake, no episode to attend")
	}
	if err := d.DiagnosePatient(p); err != nil {
		return EpisodeRecord{}, err
	}
	rec := EpisodeRecord{
		id:       newEpisodeID(),
		dateTime: time.Now(),
		patient:  p,
		attender: d,
		location: location,
		room:     p.assignedRoom, // puede ser nil si no hubo cama
	}
	d.episodes = append(d.episodes, rec)

	return rec, nil
}

// AttendSleepEmergency es la versión específica del doctor:
// usa la ubicación actual del paciente.
func (d *Doctor) AttendSleepEmergency(p *Patient) error {
	if p == nil {
		return errors.New("patient is nil")
	}
	_, err := d.Attend(p, p.currentLocation)
	return err
}

// MyEpisodes devuelve una copia, para que nadie modifique el slice interno.
func (d *Doctor) MyEpisodes() []EpisodeRecord {
	out := make([]EpisodeRecord, len(d.episodes))
	copy(out, d.episodes)
	return out
}