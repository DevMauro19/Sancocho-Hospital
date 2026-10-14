package hospital

import (
    "errors"
    "time"
)

type Doctor struct {
    Person // Embedding the Person struct
    specialty string
    patients []*Patient // Pacientes bajo su cuidado
    episodes []EpisodeRecord // numero de episodios de pacientes atendidos por el doctor
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

//Devolver la especialidad del doctor
func (d *Doctor) Specialty() string {return d.specialty}

//La funcion DiagnosePatient asigna al paciente a un doctor .
func (d *Doctor) DiagnosePatient(p *Patient) error{
    if p == nil{
        return errors.New("patient is nil")
    }
    if p.assignedDoctor != nil && p.assignedDoctor != d{
        return errors.New("patient is already assigned to another doctor")
    }
    if p.assignedDoctor == nil {
        p.assignedDoctor = d
        d.patients = append(d.patients, p)
    }
    return nil
}

//Attend implementa Attender: el objetivo de esta funcion es atender un episodio y registrarlo en el historial del doctor y del paciente.
func (d *Doctor) Attend (p* Patient, location string) (EpisodeRecord, error){
    if p == nil{
        return EpisodeRecord{}, errors.New("patient is nil")
    }
    if p.state == Awake{
        return EpisodeRecord{}, errors.New("patient is awake, nothing episode to attend")
    }
    if err:= d.DiagnosePatient(p); err != nil{
        return EpisodeRecord{}, err
    }
    rec := EpisodeRecord{
        id:        newEpisodeID(),
        dateTime:  time.Now(),
        patient:   p,
        attending: d,
        location:  location,
        assigned_room: p.assignedRoom,
    }
    d.episodes = append(d.episodes, rec)

    return rec, nil
}

//AttendSleepEmercency es la version especifica del doctor:
// usa la ubicación actual del paciente
func (d *Doctor) AttendSleepEmergency(p *Patient) (EpisodeRecord, error){
    if p == nil{
        return errors.New("patient is nil")
    }
    _, err:= d.Attend(p, p.currentLocation)
    return err
}
// MyEpisodes devuelve una copia, para que nadie modifique el slice interno.
func (d *Doctor) MyEpisodes() []EpisodeRecord {
	out := make([]EpisodeRecord, len(d.episodes))
	copy(out, d.episodes)
	return out
}