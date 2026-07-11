# lem-in — core/utils.go

## Fatal

```go
func Fatal(msg string) {
	if msg == "" {
		msg = "invalid data format"
	}
	fmt.Printf("ERROR: %s\n", msg)
	os.Exit(1)
}
```
Prints an error message and stops the program immediately.
- If no message is given, uses the default `"invalid data format"`
- Prints in the format `ERROR: <message>`
- Exits with code 1 — meaning the program ended with an error

> This function is called everywhere something goes wrong — bad input, missing rooms, invalid links etc. It keeps error handling consistent across the whole project.

---

## Atoi

```go
func Atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		Fatal(fmt.Sprintf("invalid number: '%s'", s))
	}
	return n
}
```
Converts a string to an integer safely.
- If the string is not a valid number, calls `Fatal` and stops the program
- Returns the integer if conversion succeeds

---

## IsNumber

```go
func IsNumber(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.Atoi(s)
	return err == nil
}
```
Checks if a string is a valid integer.
- Returns false for empty strings
- Returns true if the string can be converted to an integer

---

## IsComment

```go
func IsComment(line string) bool {
	return len(line) >= 1 && line[0] == '#' && (len(line) == 1 || line[1] != '#')
}
```
Checks if a line is a comment — starts with `#` but not `##`.
- A single `#` counts as a comment
- `##start` and `##end` are not comments — they start with `##`

---

## IsUnknownCommand

```go
func IsUnknownCommand(line string) bool {
	if len(line) < 2 {
		return false
	}
	if line[0] != '#' || line[1] != '#' {
		return false
	}
	if line == "##start" || line == "##end" {
		return false
	}
	return true
}
```
Checks if a line is an unknown `##` command — not `##start` or `##end`.
- Must start with `##`
- Must not be `##start` or `##end`
- Any other `##` command is unknown and should be ignored

---

## IsValidRoomName

```go
func IsValidRoomName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}
```
Checks if a room name contains only standard printable ASCII characters.
- Returns false for empty names
- Rejects any character outside the standard printable ASCII range
---