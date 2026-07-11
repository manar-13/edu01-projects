package simulate

import (
	"fmt"
	"lem-in/core"
	"lem-in/paths"
	"strings"
)

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
				CheckIfMoveEnd(d, res, ants, &sb)
			}
			res.AntNum++
		}

		if sb.Len() > 0 {
			fmt.Println(sb.String())
		}

		res.Moves++
	}
}
