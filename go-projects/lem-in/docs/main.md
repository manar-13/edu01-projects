# lem-in — main.go

## Package and Imports

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"lem-in/core"
	"lem-in/flow"
	"lem-in/parsing"
	"lem-in/simulate"
)
```
Imports the four internal packages plus standard Go tools:
- `fmt` — for printing output
- `os` — for reading command line arguments
- `path/filepath` — for checking the file extension
- `strings` — for counting dots in the filename
- `lem-in/core` — data structures and utilities
- `lem-in/flow` — BFS and flow algorithms
- `lem-in/parsing` — file reading and input parsing
- `lem-in/simulate` — ant movement simulation

---

## main

```go
if len(os.Args) < 2 {
	core.Fatal("go run main.go <input_file>")
}

path := os.Args[1]
base := filepath.Base(path)

if strings.Count(base, ".") != 1 || filepath.Ext(base) != ".txt" {
	core.Fatal("Recommended input file format: <n>.txt")
}
```
Validates the command line arguments:
- Must have exactly one argument — the input file path
- The file must have exactly one dot in its name and must end with `.txt`
- If not, print a usage error and stop

---

```go
data := &core.Data{
	Rooms:     make(map[string]*core.Room),
	Input:     parsing.ReadFile(path),
	BestSpeed: int(^uint(0) >> 1),
}
if len(data.Input) == 0 {
	core.Fatal("input file is empty")
}
```
Creates the main data structure:
- `Rooms` — empty map ready to store rooms
- `Input` — all lines from the file read immediately
- `BestSpeed` — set to the maximum possible integer value so any real result will be better

> `int(^uint(0) >> 1)` is a Go trick to get the maximum integer value without importing `math`. It inverts all bits of 0 giving all 1s, then shifts right by 1 to make it positive.

---

```go
idx := 0

parsing.ParseAnts(data, &idx)
parsing.ParseRooms(data, &idx)
parsing.ParseLinks(data, &idx)
flow.AllocateFlowPointers(data)
flow.BFSDriver(data)
flow.SetCorrectFlows(data)
for _, l := range data.Input {
	fmt.Println(l)
}
fmt.Println()
simulate.PrintMoves(data)
```
Runs the program in order:
- Parses ants, rooms, and links from the file using a shared index pointer
- Allocates flow slices for every room
- Runs BFS repeatedly to find the best set of paths
- Sets the final flow directions based on the best paths found
- Prints the original input file back to the terminal
- Prints a blank line separator
- Runs the simulation and prints each turn's ant moves

> The order is critical. Each step depends on the previous one. Parsing must finish before flow allocation, flow must finish before simulation, and the input must be printed before the moves.
---
