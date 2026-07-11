package flow

import (
	"lem-in/core"
)

func SetCorrectFlows(d *core.Data) {
	if d.BestSet == nil {
		core.Fatal("cannot set correct flows: best path set is nil")
	}

	set := d.BestSet
	i := 0

	for i < set.PathsAmount {
		if set.Lengths[i] == 0 {
			// Direct edge: start -> end, no intermediate rooms
			d.Start.Flow[i] = d.End
		} else {
			// First room after start
			d.Start.Flow[i] = set.Paths[i][0]

			// Link internal rooms
			for j := 0; j < set.Lengths[i]-1; j++ {
				set.Paths[i][j].Flow[0] = set.Paths[i][j+1]
			}

			// Last intermediate room → end
			set.Paths[i][set.Lengths[i]-1].Flow[0] = d.End
		}

		i++
	}

	// Clear remaining slots
	for i < len(d.Start.Flow) {
		d.Start.Flow[i] = nil
		i++
	}
}
