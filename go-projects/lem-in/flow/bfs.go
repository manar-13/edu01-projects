package flow

import (
	"lem-in/core"
	"lem-in/paths"
)

func resetParentsAndFlowParents(d *core.Data) {
	for _, r := range d.RoomOrder {
		r.Parent = nil
		r.FlowParent = nil
	}
}

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

func BFSDriver(d *core.Data) {
	for bfsOnce(d) {
		paths.BestPathsSetOperations(d)
	}

	if d.Start.Flow == nil || d.Start.Flow[0] == nil {
		core.Fatal("no valid flow from start to end: BFS did not find any usable paths")
	}
}
