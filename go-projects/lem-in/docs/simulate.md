# lem-in — simulate/simulate.go

## InitResDataAndAntsArr

```go
func InitResDataAndAntsArr(d *core.Data) (*core.Result, []*core.Room) {
	res := &core.Result{
		AntNum:     0,
		Moves:      0,
		Left:       0,
		Finished:   0,
		FirstPrint: true,
	}

	ants := make([]*core.Room, d.Ants)
	for i := 0; i < d.Ants; i++ {
		ants[i] = d.Start
	}

	return res, ants
}
```
Sets up the simulation state before any ants move.
- Creates a `Result` struct with all counters at zero
- Creates a slice with one entry per ant — each ant starts in the start room
- Returns both the result tracker and the ants array

---

## CheckIfMoveEnd

```go
func CheckIfMoveEnd(d *core.Data, res *core.Result, ants []*core.Room, sb *strings.Builder) {
	if ants[res.AntNum] != d.End {
		ants[res.AntNum].Occupied = true
	} else {
		res.Finished++
	}

	if res.FirstPrint {
		sb.WriteString(fmt.Sprintf("L%d-%s", res.AntNum+1, ants[res.AntNum].Name))
		res.FirstPrint = false
	} else {
		sb.WriteString(fmt.Sprintf(" L%d-%s", res.AntNum+1, ants[res.AntNum].Name))
	}
}
```
After an ant moves, checks if it reached the end and records the move for printing.
- If the ant is not at the end, marks its new room as occupied so other ants cannot enter
- If the ant reached the end, increments the finished counter
- Writes the ant's move in the format `L1-roomname` to the string builder
- The first move on a line has no leading space — subsequent moves are separated by spaces

---

## SendFromStart

```go
func SendFromStart(d *core.Data, res *core.Result, ants []*core.Room, sb *strings.Builder, pathUsedThisTurn []bool) {
	for i := 0; i < d.BestSet.PathsAmount; i++ {
		first := d.Start.Flow[i]
		if first == nil {
			continue
		}
		if pathUsedThisTurn[i] {
			continue
		}
		if !first.Occupied && paths.CanSendThisPath(d.BestSet, i, d.Ants-res.Left) {
			res.Left++
			ants[res.AntNum] = first
			pathUsedThisTurn[i] = true
			CheckIfMoveEnd(d, res, ants, sb)
			return
		}
	}
}
```
Tries to send the current ant from the start room down one of the paths.
- Loops through all paths in the best set
- Skips nil paths and paths already used this turn
- Checks that the first room of the path is not occupied and that it is worth sending an ant down this path
- If a valid path is found, moves the ant into the first room, marks the path as used this turn, and records the move
- Returns after sending one ant — each ant is sent at most once per turn

---

## PrintMoves

```go
func PrintMoves(d *core.Data) {
	res, ants := InitResDataAndAntsArr(d)

	for res.Finished != d.Ants {
		res.AntNum = 0
		res.FirstPrint = true

		var sb strings.Builder
		pathUsedThisTurn := make([]bool, d.BestSet.PathsAmount)

		for res.AntNum < d.Ants {
			if ants[res.AntNum] == d.Start {
				SendFromStart(d, res, ants, &sb, pathUsedThisTurn)
			} else if ants[res.AntNum] != d.End {
				ants[res.AntNum].Occupied = false
				ants[res.AntNum] = ants[res.AntNum].Flow[0]
				CheckIfMoveEnd(d, res, ants, sb)
			}
			res.AntNum++
		}

		if sb.Len() > 0 {
			fmt.Println(sb.String())
		}

		res.Moves++
	}
}
```
The main simulation loop — moves all ants turn by turn and prints each turn's moves.
- Initializes the simulation state
- Keeps running turns until all ants have finished
- Each turn resets the ant counter, the print flag, the string builder, and the path usage tracker
- For each ant this turn:
  - If the ant is at start — try to send it down a path
  - If the ant is in the middle — unoccupy its current room, move it to the next room, check if it finished
  - If the ant is at end — skip it, it is already done
- After processing all ants, print the turn's moves if anything moved
- Increments the move counter

> The order of processing matters. Ants already in the middle move first (implicitly — they are processed in ant number order). This ensures no two ants collide in the same room on the same turn.
---