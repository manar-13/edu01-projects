# lem-in — core/structs.go

## Room

```go
type Room struct {
	Name       string
	X, Y       int
	Links      []*Room
	Parent     *Room
	FlowParent *Room
	Flow       []*Room
	FlowFrom   *Room
	Occupied   bool
}
```
Represents one room in the ant colony.
- `Name` — the room's name from the input file
- `X, Y` — the room's coordinates (used for visualization)
- `Links` — a list of all rooms directly connected to this one by a tunnel
- `Parent` — used during BFS to track which room we came from
- `FlowParent` — used during flow path reconstruction to track reverse paths
- `Flow` — a list of rooms this room sends ants to — the chosen path direction
- `FlowFrom` — the room that sends flow into this room — used to detect and reroute existing paths
- `Occupied` — true if an ant is currently in this room

> `Flow`, `Parent`, `FlowParent`, and `FlowFrom` are all temporary fields used by the BFS and flow algorithms. They have no meaning from the input file — they are computed during path finding.

---

## PathsSet

```go
type PathsSet struct {
	Paths       [][]*Room
	Lengths     []int
	PathsAmount int
}
```
Represents one complete set of paths found by BFS.
- `Paths` — a 2D list where each row is one path (a list of rooms from start to end)
- `Lengths` — the length of each path (number of intermediate rooms, not counting start and end)
- `PathsAmount` — how many paths are in this set

---

## Result

```go
type Result struct {
	Finished   int
	AntNum     int
	Moves      int
	Left       int
	FirstPrint bool
}
```
Tracks the state of the simulation while ants are moving.
- `Finished` — how many ants have reached the end room
- `AntNum` — the index of the ant currently being processed this turn
- `Moves` — how many turns have passed
- `Left` — how many ants have left the start room so far
- `FirstPrint` — true if this is the first move printed on the current line

---

## Data

```go
type Data struct {
	Ants      int
	Rooms     map[string]*Room
	Input     []string
	Start     *Room
	End       *Room
	BestSet   *PathsSet
	BestSpeed int
	RoomOrder []*Room
}
```
The main data structure passed through the entire program.
- `Ants` — the number of ants from the input file
- `Rooms` — a map of all rooms keyed by name for fast lookup
- `Input` — the raw lines from the input file
- `Start` — pointer to the start room
- `End` — pointer to the end room
- `BestSet` — the best combination of paths found so far
- `BestSpeed` — the number of moves the best set takes — lower is better
- `RoomOrder` — rooms in the order they were read from the file

> `Data` is created in `main.go` and passed to every function. This avoids global variables and keeps the code clean.
---
