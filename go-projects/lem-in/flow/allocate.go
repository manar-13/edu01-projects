package flow

import (
	"lem-in/core"
)

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

func allocateOtherFlows(room *core.Room) {
	room.Flow = make([]*core.Room, 2)
}

func AllocateFlowPointers(d *core.Data) {
	for _, r := range d.RoomOrder {
		if r == d.Start {
			allocateStartFlows(d, r)
		} else {
			allocateOtherFlows(r)
		}
	}
}
