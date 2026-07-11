# lem-in — flow/bfs.go

## resetParentsAndFlowParents

```go
func resetParentsAndFlowParents(d *core.Data) {
	for _, r := range d.RoomOrder {
		r.Parent = nil
		r.FlowParent = nil
	}
}
```
Clears the `Parent` and `FlowParent` fields of every room before each BFS run.
- BFS uses these fields to track how it reached each room
- Resetting them before each run makes sure old data from the previous BFS does not interfere

---

## iterateLinks

```go
func iterateLinks(r *core.Room, queue *[]*core.Room) {
	for _, link := range r.Links {
		if positiveFlow(r.Flow, link) {
			continue
		} else if r.FlowFrom != nil && r.FlowParent == nil {
			foundOldPath(queue, r)
			return
		} else if r.FlowFrom != nil && r.FlowParent != nil {
			canGoEverywhere(r, link, queue)
		} else if link.Parent == nil {
			visitUsingUnusedEdge(queue, r, link)
		}
	}
}
```
Processes all neighbors of the current room during BFS.
- Skips any neighbor that is already in the current room's flow — that edge is taken
- If the current room has a `FlowFrom` but no `FlowParent` — we found an existing path to reverse, call `foundOldPath` and stop
- If the current room has both `FlowFrom` and `FlowParent` — we are on a reverse path, try to go anywhere with `canGoEverywhere`
- Otherwise use a normal unused edge with `visitUsingUnusedEdge`

> This logic implements a modified BFS that can find augmenting paths through existing flow — the core idea behind the max-flow algorithm used here.

---

## bfsOnce

```go
func bfsOnce(d *core.Data) bool {
	if d.Start == nil || d.End == nil ||
		len(d.Start.Links) == 0 || len(d.End.Links) == 0 ||
		d.Start == d.End {
		core.Fatal("invalid start or end room configuration")
	}

	resetParentsAndFlowParents(d)
	d.End.Parent = nil

	queue := []*core.Room{d.Start}

	for len(queue) > 0 && d.End.Parent == nil {
		current := queue[0]
		queue = queue[1:]
		iterateLinks(current, &queue)
	}

	if d.End.Parent == nil {
		return false
	}

	setFlows(d)
	return true
}
```
Runs one complete BFS from start to end and updates the flow if a path is found.
- Validates that start and end rooms exist, have links, and are not the same room
- Resets all parent pointers
- Starts BFS from the start room
- Processes rooms one by one until the end room is reached or the queue is empty
- If the end room has no parent after BFS, no path was found — return false
- If a path was found, call `setFlows` to record it and return true

---

## BFSDriver

```go
func BFSDriver(d *core.Data) {
	for bfsOnce(d) {
		paths.BestPathsSetOperations(d)
	}

	if d.Start.Flow == nil || d.Start.Flow[0] == nil {
		core.Fatal("no valid flow from start to end: BFS did not find any usable paths")
	}
}
```
Runs BFS repeatedly until no more paths can be found.
- Each call to `bfsOnce` finds one new augmenting path and adds it to the flow
- After each successful BFS, `BestPathsSetOperations` checks if the current set of paths is the best one found so far
- When BFS can no longer find a path, the loop ends
- If no paths were found at all, stop with an error

> Running BFS multiple times is how the program finds multiple paths. Each run adds one more path to the flow network. After all runs, the best combination of paths is selected based on how fast it can move all the ants.
---
