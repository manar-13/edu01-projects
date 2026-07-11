# lem-in — flow/allocate.go

## allocateStartFlows

```go
func allocateStartFlows(d *core.Data, room *core.Room) {
	if d.End == nil {
		core.Fatal("cannot allocate start flows: end room is nil")
	}
	size := len(d.End.Links)
	if len(room.Links) > size {
		size = len(room.Links)
	}
	room.Flow = make([]*core.Room, size+1)
}
```
Allocates the `Flow` slice for the start room.
- The start room needs a larger `Flow` slice because it can send ants down multiple paths at once
- The size is the larger of the start room's links count or the end room's links count, plus one
- This gives enough slots for all possible outgoing paths

---

## allocateOtherFlows

```go
func allocateOtherFlows(room *core.Room) {
	room.Flow = make([]*core.Room, 2)
}
```
Allocates the `Flow` slice for every non-start room.
- Non-start rooms only need 2 slots — one for the forward flow and one reserved slot
- This keeps memory usage minimal for regular rooms

---

## AllocateFlowPointers

```go
func AllocateFlowPointers(d *core.Data) {
	for _, r := range d.RoomOrder {
		if r == d.Start {
			allocateStartFlows(d, r)
		} else {
			allocateOtherFlows(r)
		}
	}
}
```
Allocates `Flow` slices for every room before BFS runs.
- Loops through all rooms in the order they were read
- Gives the start room a larger slice and every other room a smaller one

> This must run before BFS because BFS reads and writes to the `Flow` slices of every room. Without allocating them first the program would crash.
---
