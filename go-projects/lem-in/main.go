package main

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

func main() {
	if len(os.Args) < 2 {
		core.Fatal("go run main.go <input_file>")
	}

	path := os.Args[1]
	base := filepath.Base(path)

	if strings.Count(base, ".") != 1 || filepath.Ext(base) != ".txt" {
		core.Fatal("Recommended input file format: <name>.txt")
	}

	data := &core.Data{
		Rooms:     make(map[string]*core.Room),
		Input:     parsing.ReadFile(path),
		BestSpeed: int(^uint(0) >> 1),
	}
	if len(data.Input) == 0 {
		core.Fatal("input file is empty")
	}

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
}
