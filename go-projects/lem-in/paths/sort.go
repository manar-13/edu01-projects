package paths

import (
	"lem-in/core"
)

func SortPathsShortToLong(set *core.PathsSet) {
	n := set.PathsAmount

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if set.Lengths[j] < set.Lengths[i] {
				set.Lengths[i], set.Lengths[j] = set.Lengths[j], set.Lengths[i]
				set.Paths[i], set.Paths[j] = set.Paths[j], set.Paths[i]
				i = -1
				break
			}
		}
	}
}
