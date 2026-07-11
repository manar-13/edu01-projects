package paths

import (
	"lem-in/core"
)

func SaveCurrentPathsSet(d *core.Data, set *core.PathsSet) {
	i := 0
	for _, head := range d.Start.Flow {
		if head == nil {
			continue
		}
		j := 0
		room := head
		for room != d.End {
			set.Paths[i][j] = room
			room = room.Flow[0]
			j++
		}
		i++
	}
}
