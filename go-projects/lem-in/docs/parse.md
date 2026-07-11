# lem-in — parsing/parse.go

## ParseAnts

```go
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

	if n > 100000 {
		core.Fatal(fmt.Sprintf(
			"ant count too high: %d (maximum allowed is 100000)", n,
		))
	}

	d.Ants = n
	*idx++
}
```
Reads the number of ants from the first non-comment line.
- Skips any comment or unknown command lines at the start
- If we run out of lines without finding the ants number, stop with an error
- Checks the line is a valid number — if not, stop with an error
- Checks the number is positive and not above 100000
- Stores the ant count in `data.Ants` and moves the index forward

> `idx` is a pointer to an integer that tracks which line we are currently reading. Passing it as a pointer means `ParseAnts`, `ParseRooms`, and `ParseLinks` all share the same position in the file and pick up where the previous one left off.

---

## getOrCreateRoom

```go
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
```
Gets an existing room by name or creates a new one.
- Validates the room name — must be printable ASCII
- Room names cannot start with `L` (reserved for ant labels like `L1`) or `#`
- Room names cannot contain `-` (reserved for link definitions)
- If the room already exists in the map, return it
- Otherwise create a new room, store its coordinates, add it to the map and order list

---

## parseRoomLine

```go
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
		core.Fatal(fmt.Sprintf("room coordinates cannot be negative: (%d, %d)", x, y))
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
```
Parses one room line in the format `name x y`.
- Splits the line into 3 parts — if not exactly 3, stop with an error
- Validates that x and y are numbers and not negative
- Checks for duplicate room names
- Checks for duplicate coordinates — two rooms cannot share the same position
- Creates and returns the room

---

## ParseRooms

```go
func ParseRooms(d *core.Data, idx *int) {
	lines := d.Input
	var pendingStart, pendingEnd bool

	for *idx < len(lines) {
		line := lines[*idx]

		if line == "" || core.IsComment(line) || core.IsUnknownCommand(line) {
			*idx++
			continue
		}

		if pendingStart || pendingEnd {
			r := parseRoomLine(d, line)
			if pendingStart {
				d.Start = r
			} else {
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

		break
	}

	if d.Start == nil || d.End == nil {
		core.Fatal("missing ##start or ##end room definition")
	}
}
```
Reads all room definitions from the file.
- Skips empty lines, comments, and unknown commands
- When `##start` is seen, sets `pendingStart = true` — the next room line becomes the start room
- When `##end` is seen, sets `pendingEnd = true` — the next room line becomes the end room
- Any line containing a space is treated as a room definition `name x y`
- When a line with no space is found, the rooms section is over — break and let `ParseLinks` take over
- After the loop, checks that both start and end rooms were found

---

## addLink

```go
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
```
Adds a one-directional link from room `a` to room `b`.
- First checks if the link already exists — if it does, stop with an error
- Adds `b` to `a`'s list of linked rooms

> `ParseLinks` calls `addLink(r1, r2)` and `addLink(r2, r1)` to make the link bidirectional.

---

## ParseLinks

```go
func ParseLinks(d *core.Data, idx *int) {
	lines := d.Input
	for *idx < len(lines) {
		line := lines[*idx]

		if line == "" || core.IsComment(line) || core.IsUnknownCommand(line) {
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
```
Reads all link definitions from the file in the format `room1-room2`.
- Skips empty lines, comments, and unknown commands
- If `##start` or `##end` appear here, stop with an error — they belong in the rooms section
- Every line must contain exactly one `-` separating two room names
- Both rooms must already exist in the rooms map
- A room cannot link to itself
- Adds the link in both directions
---
