# lem-in — flow/set_correct.go

## SetCorrectFlows

```go
func SetCorrectFlows(d *core.Data) {
	if d.BestSet == nil {
		core.Fatal("cannot set correct flows: best path set is nil")
	}

	set := d.BestSet
	i := 0

	for i < set.PathsAmount {
		if set.Lengths[i] == 0 {
			d.Start.Flow[i] = d.End
		} else {
			d.Start.Flow[i] = set.Paths[i][0]

			for j := 0; j < set.Lengths[i]-1; j++ {
				set.Paths[i][j].Flow[0] = set.Paths[i][j+1]
			}

			set.Paths[i][set.Lengths[i]-1].Flow[0] = d.End
		}
		i++
	}

	for i < len(d.Start.Flow) {
		d.Start.Flow[i] = nil
		i++
	}
}
```
After all BFS runs are done, this sets the final flow directions based on the best path set found.
- If `BestSet` is nil, stop with an error
- For each path in the best set:
  - If the path has length 0 (start connects directly to end), set `d.Start.Flow[i] = d.End`
  - Otherwise set the start room's flow to the first room of the path
  - Chain each intermediate room's flow to the next room in the path
  - Set the last intermediate room's flow to the end room
- Clear any remaining flow slots in the start room that are not used by the best set

> This is the final step before simulation. It rewrites all the flow directions to match the best path combination, overwriting anything BFS may have left from exploring non-optimal combinations.
---
