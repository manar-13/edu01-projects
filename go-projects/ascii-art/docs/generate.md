# ascii-art — art/generate.go

## GenerateASCIIArt

```go
func GenerateASCIIArt(input string, font map[rune][]string) string {
	var result strings.Builder
	lines := strings.Split(input, "\\n")
```
Takes the input text and the font map and creates an empty builder to collect the output. Splits the input by `\n` to handle multiple lines.

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
Loop through each line. If the line is empty and is not the first line, add a new line to the output and skip to the next one.

> This handles cases like `"Hello\n\nWorld"` where there is an empty line between two words — it adds that empty line to the output.

---

```go
		for h := 0; h < 8; h++ {
			for _, ch := range line {
				if art, ok := font[ch]; ok {
					result.WriteString(art[h])
				} else {
					result.WriteString(strings.Repeat(" ", 8))
				}
			}
			result.WriteString("\n")
		}
	}
	return result.String()
}
```
For each line of text, loop through all 8 rows of the ASCII art:
- For each character in the line, look it up in the font map
- If found, add that character's art for the current row
- If not found, add 8 spaces as a placeholder
- After all characters in a row are done, add a new line
- When all 8 rows are done, move to the next line of text
- Return the complete ASCII art as one string

> The outer loop goes row by row (0 to 7). The inner loop goes character by character. This is how all characters end up side by side — we print one row of every character before moving to the next row.

---

## PrintASCIIArt

```go
func PrintASCIIArt(input string, font map[rune][]string) {
	fmt.Print(GenerateASCIIArt(input, font))
}
```
A simple wrapper around `GenerateASCIIArt` that prints the result directly to the terminal instead of returning it as a string.
---