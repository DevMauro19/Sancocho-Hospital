package hospital

import (
	"fmt"
	"sync/atomic"
	"time"
)

type EpisodeRecord struct {
	id       string
	dateTime time.Time
	patient  *Patient
	attender Attender
	location string
	room     *Room // puede ser nil
}

func (e EpisodeRecord) Summary() string {
	room := "no room"
	if e.room != nil {
		room = fmt.Sprintf("room %d", e.room.Number())
	}
	pid := "?"
	if e.patient != nil {
		pid = e.patient.ID()
	}
	return fmt.Sprintf("%s  %s  %s  -> %s", e.dateTime.Format("2006-01-02 15:04"), pid, e.location, room)
}

var episodeCounter atomic.Int64

func newEpisodeID() string {
	return fmt.Sprintf("E-%03d", episodeCounter.Add(1))
}