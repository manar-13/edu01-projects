# lem-in — paths/lengths.go

## CalculatePathsLengths

```go
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
```
Calculates the length of each path — the number of intermediate rooms not counting start and end.
- Loops through each starting point in `d.Start.Flow`
- If the path goes directly to end (length 0), records 0
- Otherwise follows the `Flow` pointers room by room until the end room is reached, counting steps
- Stores the length in `set.Lengths[i]`

---

## AllocatePathsArrays

```go
func AllocatePathsArrays(set *core.PathsSet) {
	for i := 0; i < set.PathsAmount; i++ {
		set.Paths[i] = make([]*core.Room, set.Lengths[i])
	}
}
```
Allocates the inner arrays for each path based on its calculated length.
- For each path, creates a slice of exactly the right size to hold its intermediate rooms
- Must run after `CalculatePathsLengths` so the lengths are known
---
