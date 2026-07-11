package paths

import (
	"lem-in/core"
)

func CalculatePathsLengths(d *core.Data, set *core.PathsSet) {
	i := 0
	for _, head := range d.Start.Flow {
		if head == nil {
			continue
		}

		if head == d.End {
			set.Lengths[i] = 0
			i++
			continue
		}

		length := 1
		iterator := head
		for iterator.Flow[0] != d.End {
			iterator = iterator.Flow[0]
			length++
		}
		set.Lengths[i] = length
		i++
	}
}

func AllocatePathsArrays(set *core.PathsSet) {
	for i := 0; i < set.PathsAmount; i++ {
		set.Paths[i] = make([]*core.Room, set.Lengths[i])
	}
}
