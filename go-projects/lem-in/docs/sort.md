# lem-in — paths/sort.go

## SortPathsShortToLong

```go
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
```
Sorts the paths from shortest to longest using a bubble sort style algorithm.
- Compares every pair of paths
- If a later path is shorter than an earlier one, swaps both the lengths and the path arrays
- Resets `i` to `-1` after a swap so the outer loop restarts from 0 — this ensures the sort is stable and complete
- Stops when no more swaps are needed

> Sorting paths shortest to longest is important for `CanSendThisPath` and `CheckIfCurrentIsBest`. The ant distribution logic assumes shorter paths come first so it can correctly calculate how many ants each path should receive.
---
