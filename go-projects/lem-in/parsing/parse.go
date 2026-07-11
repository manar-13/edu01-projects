package parsing

import (
	"fmt"
	"strings"

	"lem-in/core"
)

func ParseAnts(d *core.Data, idx *int) {
	lines := d.Input
	for *idx < len(lines) &&
		(core.IsComment(lines[*idx]) || core.IsUnknownCommand(lines[*idx])) {
		*idx++
	}
	if *idx >= len(lines) {
		core.Fatal("missing ants line")
	}

	line := strings.TrimSpace(lines[*idx])
	if !core.IsNumber(line) {
		core.Fatal("number of ants is not numeric")
	}

	n := core.Atoi(line)
	if n <= 0 {
		core.Fatal("number of ants must be positive")
	}

	if n <= 0 {
		core.Fatal("number of ants must be positive")
	}

	if n > 100000 {
		core.Fatal(fmt.Sprintf(
			"ant count too high: %d (maximum allowed is 100000)",
			n,
		))
	}

	d.Ants = n
	*idx++
}

func getOrCreateRoom(d *core.Data, name string, x, y *int) *core.Room {

	if !core.IsValidRoomName(name) {
		core.Fatal("invalid room name")
	}

	if name[0] == 'L' || name[0] == '#' {
		core.Fatal("room name cannot start with 'L' or '#'")
	}
	if strings.Contains(name, "-") {
		core.Fatal("room name cannot contain '-'")
	}

	if r, ok := d.Rooms[name]; ok {
		return r
	}

	r := &core.Room{Name: name}
	if x != nil {
		r.X = *x
	}
	if y != nil {
		r.Y = *y
	}

	d.Rooms[name] = r
	d.RoomOrder = append(d.RoomOrder, r)
	return r
}

func parseRoomLine(d *core.Data, line string) *core.Room {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		core.Fatal("invalid room line format: expected 'name x y'")
	}
	if !core.IsNumber(parts[1]) || !core.IsNumber(parts[2]) {
		core.Fatal("room coordinates must be numeric")
	}

	name := parts[0]
	x := core.Atoi(parts[1])
	y := core.Atoi(parts[2])
	if x < 0 || y < 0 {
		core.Fatal(fmt.Sprintf(
			"room coordinates cannot be negative: (%d, %d)",
			x, y,
		))
	}

	if _, exists := d.Rooms[name]; exists {
		core.Fatal("duplicate room name: " + name)
	}

	for _, existing := range d.Rooms {
		if existing.X == x && existing.Y == y {
			core.Fatal(fmt.Sprintf("duplicate coordinates (%d,%d)", x, y))
		}
	}

	return getOrCreateRoom(d, name, &x, &y)
}

func ParseRooms(d *core.Data, idx *int) {
	lines := d.Input
	var pendingStart, pendingEnd bool

	for *idx < len(lines) {
		line := lines[*idx]

		if line == "" {
			*idx++
			continue
		}
		if core.IsComment(line) || core.IsUnknownCommand(line) {
			*idx++
			continue
		}

		if pendingStart || pendingEnd {
			if strings.HasPrefix(line, "#") || !strings.Contains(line, " ") {
				core.Fatal("expected room after ##start or ##end")
			}
			r := parseRoomLine(d, line)
			if pendingStart {
				if d.Start != nil {
					core.Fatal("multiple ##start declarations")
				}
				d.Start = r
			} else {
				if d.End != nil {
					core.Fatal("multiple ##end declarations")
				}
				d.End = r
			}
			pendingStart, pendingEnd = false, false
			*idx++
			continue
		}

		if line == "##start" {
			pendingStart = true
			*idx++
			continue
		}
		if line == "##end" {
			pendingEnd = true
			*idx++
			continue
		}

		if strings.Contains(line, " ") {
			parseRoomLine(d, line)
			*idx++
			continue
		}

		break // links begin
	}

	if d.Start == nil || d.End == nil {
		core.Fatal("missing ##start or ##end room definition")
	}
}

func addLink(a, b *core.Room) {
	for _, x := range a.Links {
		if x == b {
			core.Fatal(fmt.Sprintf(
				"duplicate link detected: room '%s' is already connected to '%s'",
				a.Name, b.Name,
			))
		}
	}
	a.Links = append(a.Links, b)
}

func ParseLinks(d *core.Data, idx *int) {
	lines := d.Input
	if *idx >= len(lines) {
		core.Fatal("no links section found")
	}

	for *idx < len(lines) {
		line := lines[*idx]

		if line == "" {
			*idx++
			continue
		}
		if core.IsComment(line) || core.IsUnknownCommand(line) {
			*idx++
			continue
		}
		if line == "##start" || line == "##end" {
			core.Fatal("##start or ##end cannot appear inside link definitions")
		}
		if !strings.Contains(line, "-") {
			core.Fatal("invalid link format: expected 'room1-room2'")
		}

		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			core.Fatal("invalid link: too many '-' characters")
		}

		r1 := d.Rooms[parts[0]]
		r2 := d.Rooms[parts[1]]
		if r1 == nil || r2 == nil {
			core.Fatal("link references undefined room")
		}
		if r1 == r2 {
			core.Fatal("room cannot link to itself")
		}

		addLink(r1, r2)
		addLink(r2, r1)

		*idx++
	}
}
