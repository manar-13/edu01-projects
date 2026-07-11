package flow

import (
	"lem-in/core"
)

func setFlows(d *core.Data) {
	current := d.End
	parent := current.Parent

	for current != d.Start {
		if parent == d.Start {
			i := 0
			for i < len(d.Start.Flow) && d.Start.Flow[i] != nil {
				i++
			}
			if i < len(d.Start.Flow) {
				d.Start.Flow[i] = current
			}
		}

		if current.FlowParent != nil && current.FlowParent.FlowFrom == current {
			current.FlowParent.FlowFrom = nil
			current = current.FlowParent
			parent = current.Parent
		} else {
			if parent != d.Start {
				parent.Flow[0] = current
			}
			if current != d.End {
				current.FlowFrom = parent
			}
			current = parent
			if current != nil {
				parent = current.Parent
			} else {
				break
			}
		}
	}
}
