package flow

import (
	"lem-in/core"
)

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

func foundOldPath(queue *[]*core.Room, r *core.Room) {
	if r.FlowFrom == nil {
		return
	}
	*queue = append(*queue, r.FlowFrom)
	r.FlowFrom.FlowParent = r
}

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
