# lem-in — paths/paths.go

## BuildPathsSetStructure

```go
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
```
Creates an empty `PathsSet` structure based on how many paths currently exist in the flow.
- Counts how many non-nil entries are in `d.Start.Flow` — each one is the start of a path
- If no paths exist, returns an empty set
- Allocates the `Paths` and `Lengths` slices with the right size
- Returns the empty structure ready to be filled by other functions

> This is always the first step when saving a path set. It just creates the container — the actual paths are saved by `SaveCurrentPathsSet` and the lengths are calculated by `CalculatePathsLengths`.
---
