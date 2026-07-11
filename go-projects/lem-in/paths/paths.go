package paths

import (
	"lem-in/core"
)

func BuildPathsSetStructure(d *core.Data) *core.PathsSet {
	set := &core.PathsSet{}
	for _, head := range d.Start.Flow {
		if head != nil {
			set.PathsAmount++
		}
	}
	if set.PathsAmount == 0 {
		return set
	}

	set.Paths = make([][]*core.Room, set.PathsAmount)
	set.Lengths = make([]int, set.PathsAmount)

	return set
}
