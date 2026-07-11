# lem-in — flow/set_flow.go

## setFlows

```go
func setFlows(d *core.Data) {
	current := d.End
	parent := current.Parent

	for current != d.Start {
		if parent == d.Start {
			i := 0
			for i < len(d.Start.Flow) && d.Start.Flow[i] != nil {
				i++
			}
			if i < len(d.Start.Flow) {
				d.Start.Flow[i] = current
			}
		}

		if current.FlowParent != nil && current.FlowParent.FlowFrom == current {
			current.FlowParent.FlowFrom = nil
			current = current.FlowParent
			parent = current.Parent
		} else {
			if parent != d.Start {
				parent.Flow[0] = current
			}
			if current != d.End {
				current.FlowFrom = parent
			}
			current = parent
			if current != nil {
				parent = current.Parent
			} else {
				break
			}
		}
	}
}
```
Traces back the path found by BFS and records the flow direction for each room.
- Starts at the end room and walks backwards using `Parent` pointers
- When the parent is the start room, finds the next free slot in `d.Start.Flow` and records the first room of this path
- If the current room has a `FlowParent` that points back to it (a reverse flow situation), clears the reverse connection and follows the flow parent instead
- Otherwise records the forward flow direction — `parent.Flow[0] = current` means "from this room, go to current"
- Also records `FlowFrom` on each room so we know where flow enters from

> This is called after every successful BFS run. It converts the BFS parent chain into actual flow directions stored on each room, which are then used by `SetCorrectFlows` and the simulator.
---
