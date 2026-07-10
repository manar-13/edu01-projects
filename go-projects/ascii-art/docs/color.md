# ascii-art — art/color.go

## ColorMap

```go
var ColorMap = map[string]string{
	"black":         "\033[30m",
	"red":           "\033[31m",
	"green":         "\033[32m",
	"yellow":        "\033[33m",
	"blue":          "\033[34m",
	"magenta":       "\033[35m",
	"cyan":          "\033[36m",
	"white":         "\033[37m",
	"gray":          "\033[90m",
	"brightred":     "\033[91m",
	"brightgreen":   "\033[92m",
	"brightyellow":  "\033[93m",
	"brightblue":    "\033[94m",
	"brightmagenta": "\033[95m",
	"brightcyan":    "\033[96m",
	"brightwhite":   "\033[97m",
	"orange":        "\033[38;5;208m",
	"rose":          "\033[38;5;205m",
	"sky":           "\033[38;5;117m",
	"lime":          "\033[38;5;154m",
	"gold":          "\033[38;5;220m",
}
```
A map that connects color names to their ANSI color codes.
- The key is the color name like `"red"`
- The value is the ANSI escape code like `"\033[31m"` which tells the terminal to switch to that color

> ANSI escape codes are special sequences that terminals understand. `\033[31m` means "start printing in red." `\033[0m` means "reset back to normal color."

---

## Colorize

```go
func Colorize(segment, color string) string {
	colorCode, ok := ColorMap[strings.ToLower(color)]
	if !ok {
		return segment
	}
	return fmt.Sprintf("%s%s\033[0m", colorCode, segment)
}
```
Wraps a piece of text with a color code and a reset code.
- Looks up the color name in `ColorMap`
- If the color is not found, return the text as is without any color
- If found, wrap the text like this: `COLOR_CODE + text + RESET_CODE`

> The reset code `\033[0m` at the end is important — without it, everything printed after this text would also be colored.

---

## FindSubStringIndices

```go
func FindSubStringIndices(str, substr string) []int {
	var indices []int
	runes := []rune(str)
	subRunes := []rune(substr)
	for i := 0; i <= len(runes)-len(subRunes); i++ {
		match := true
		for j := range subRunes {
			if runes[i+j] != subRunes[j] {
				match = false
				break
			}
		}
		if match {
			indices = append(indices, i)
		}
	}
	return indices
}
```
Finds all the positions where the substring appears inside the full text.
- Converts both the text and substring into runes for safe character handling
- Slides through the text character by character
- At each position, checks if the next characters match the substring
- If they match, saves that position
- Returns a list of all starting positions where the substring was found

> For example if the text is `"a king kitten have kit"` and the substring is `"kit"`, this returns `[7, 14, 19]` — the positions where `kit` starts.

---

## isInRange

```go
func isInRange(index int, indices []int, length int) bool {
	for _, start := range indices {
		if index >= start && index < start+length {
			return true
		}
	}
	return false
}
```
Checks if a character at a given position falls inside any of the substring matches.
- Takes a character index, the list of match positions, and the substring length
- For each match position, checks if the index falls within that match
- Returns true if the character is part of a match, false if not

> For example if `kit` starts at position 7 and has length 3, then positions 7, 8, and 9 are all inside the match. This function checks that.

---

## DrawColorASCII

```go
func DrawColorASCII(input, color, substr string, font map[rune][]string) string {
	var result strings.Builder
	lines := strings.Split(input, "\\n")
```
The main color function. Splits the input by `\n` to handle multiple lines and creates an empty builder to collect the output.

---

```go
	for i, line := range lines {
		if line == "" {
			if i > 0 {
				result.WriteString("\n")
			}
			continue
		}
```
Loop through each line. If the line is empty and is not the first line, add a new line to the output and skip to the next line.

---

```go
		indices := FindSubStringIndices(line, substr)
		runes := []rune(line)
```
Find all the positions where the substring appears in this line. Convert the line to runes for safe character by character access.

---

```go
		for row := 0; row < 8; row++ {
			for idx, ch := range runes {
				art, ok := font[ch]
				if !ok {
					continue
				}
				segment := art[row]
				shouldColor := substr == "" || isInRange(idx, indices, len([]rune(substr)))
				if shouldColor && color != "" && ch != ' ' {
					result.WriteString(Colorize(segment, color))
				} else {
					result.WriteString(segment)
				}
			}
			result.WriteString("\n")
		}
	}
	return result.String()
}
```
For each of the 8 rows, loop through every character in the line:
- Look up the ASCII art for that character
- Get the current row of that art
- Decide if this character should be colored:
  - If no substring was given — color everything
  - If a substring was given — only color characters that fall inside a match
- Spaces are never colored even if they are inside the substring
- Write the colored or normal segment to the result
- After all 8 rows are done, return the complete colored ASCII art
---