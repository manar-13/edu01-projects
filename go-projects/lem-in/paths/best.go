package paths

import (
	"lem-in/core"
)

func CanSendThisPath(set *core.PathsSet, i int, ants int) bool {
	if i == 0 {
		return true
	}

	sum := 0
	for j := i - 1; j >= 0; j-- {
		if set.Lengths[j] < set.Lengths[i] {
			sum += set.Lengths[i] - set.Lengths[j]
		}
	}

	return ants > sum
}

func checkLongestMove(d *core.Data, set *core.PathsSet, antsToPath []int) {
	longest := 0

	for i := 0; i < set.PathsAmount; i++ {
		val := antsToPath[i] + set.Lengths[i]
		if val > longest {
			longest = val
		}
	}

	if d.BestSet == nil || longest < d.BestSpeed {
		d.BestSet = set
		d.BestSpeed = longest
	}
}

func CheckIfCurrentIsBest(d *core.Data, set *core.PathsSet) {
	ants := d.Ants
	antsToPath := make([]int, set.PathsAmount)

	for ants > 0 {
		for i := 0; i < set.PathsAmount && ants > 0; i++ {
			if CanSendThisPath(set, i, ants) {
				antsToPath[i]++
				ants--
			}
		}
	}

	checkLongestMove(d, set, antsToPath)
}

func BestPathsSetOperations(d *core.Data) {
	set := BuildPathsSetStructure(d)
	if set.PathsAmount == 0 {
		return
	}

	CalculatePathsLengths(d, set)
	AllocatePathsArrays(set)
	SaveCurrentPathsSet(d, set)
	SortPathsShortToLong(set)
	CheckIfCurrentIsBest(d, set)
}
