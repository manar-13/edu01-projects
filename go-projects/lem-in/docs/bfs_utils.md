# lem-in — flow/bfs_utils.go

## positiveFlow

```go
func positiveFlow(flows []*core.Room, link *core.Room) bool {
	for _, f := range flows {
		if f == nil {
			break
		}
		if f == link {
			return true
		}
	}
	return false
}
```
Checks if a room is already in the current flow list.
- Loops through the flow slice until a nil entry is found
- Returns true if the given room is already in the flow — meaning this edge is already being used by a path
- Used in BFS to avoid reusing edges that are already part of a path

---

## foundOldPath

```go
func foundOldPath(queue *[]*core.Room, r *core.Room) {
	if r.FlowFrom == nil {
		return
	}
	*queue = append(*queue, r.FlowFrom)
	r.FlowFrom.FlowParent = r
}
```
When BFS reaches a room that is on an existing path, it reverses that path edge.
- Adds the room that was sending flow into this room back onto the BFS queue
- Sets `FlowParent` to track the reverse direction
- This is part of the augmenting path technique — finding a new path by reversing parts of existing ones

---

## canGoEverywhere

```go
func canGoEverywhere(current, link *core.Room, queue *[]*core.Room) {
	if link.Parent != nil {
		return
	}
	*queue = append(*queue, link)
	if current.FlowFrom == link {
		link.FlowParent = current
	} else {
		link.Parent = current
	}
}
```
Adds a neighboring room to the BFS queue when the current room has both a flow parent and a flow from.
- Skips rooms that have already been visited
- If the neighbor is the room that sends flow into the current room, sets `FlowParent` to track the reverse
- Otherwise sets the normal `Parent` for forward traversal

---

## visitUsingUnusedEdge

```go
func visitUsingUnusedEdge(queue *[]*core.Room, current, link *core.Room) {
	if len(current.Flow) > 0 && current.Flow[0] == link {
		return
	}
	if link == current.Parent {
		return
	}
	*queue = append(*queue, link)
	link.Parent = current
}
```
Adds a neighboring room to the BFS queue using a normal unused edge.
- Skips the neighbor if it is already the direction the current room sends flow — to avoid loops
- Skips if it is the parent we came from
- Otherwise adds it to the queue and records the current room as its parent
---
