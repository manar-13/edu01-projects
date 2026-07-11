# lem-in — paths/save.go

## SaveCurrentPathsSet

```go
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
```
Saves the intermediate rooms of each current path into the `PathsSet`.
- Loops through each path starting from `d.Start.Flow`
- For each path, follows the `Flow` pointers room by room until the end room
- Saves each intermediate room into `set.Paths[i][j]`
- Does not save the start or end rooms — only the rooms in between

> This creates a snapshot of the current paths. It is important because BFS will modify the flow on the next run. Saving the paths lets us compare the current set against future sets and keep the best one.
---
