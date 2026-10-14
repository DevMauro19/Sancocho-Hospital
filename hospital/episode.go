package hospital

import (
    "fmt"
    "sync/atomic"
)
//Encargado de los contadores de episodios, para generar IDs únicos para cada episodio. Se utiliza un contador atómico para garantizar la seguridad.
var episodeCounter atomic.Int64 // import "sync/atomic" y "fmt"

func newEpisodeID() string {
	return fmt.Sprintf("E-%03d", episodeCounter.Add(1))
}