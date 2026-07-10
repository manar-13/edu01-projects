# ascii-art-web — art_web/printer.go

## GenerateASCIIArt

```go
func GenerateASCIIArt(input string, font map[rune][]string) (string, error) {
	var result strings.Builder

	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")

	lines := strings.Split(input, "\n")
```
Takes the input text and the font map and creates an empty builder to collect the output.
- Replaces `\r\n` (Windows line endings) and `\r` (old Mac line endings) with `\n` — this makes sure line breaks work the same on all operating systems
- Splits the input into separate lines wherever there is a `\n`

> We clean the line endings first because web browsers sometimes send `\r\n` instead of just `\n` depending on the operating system.

---

```go
	for lineNum, line := range lines {
		if lineNum > 0 {
			result.WriteString("\n")
		}

		if len(line) == 0 {
			continue
		}
```
Loop through every line of the input:
- If this is not the first line, add a new line before it to separate the blocks
- If the line is empty, skip it — nothing to draw

---

```go
		for layer := 0; layer < 8; layer++ {
			for _, char := range line {
				if art, ok := font[char]; ok && layer < len(art) {
					result.WriteString(art[layer])
				} else {
					return "", fmt.Errorf("character %q not found in font", char)
				}
			}
			result.WriteString("\n")
		}
	}
	return result.String(), nil
}
```
For each line of text, loop through all 8 rows of the ASCII art:
- The outer loop goes through rows 0 to 7
- The inner loop goes through every character in the line
- For each character, look it up in the font map and add its current row to the result
- If a character is not found in the font map, return an error immediately
- After all characters in a row are done, add a new line
- When all 8 rows are done, move to the next line of text
- Return the complete ASCII art as one string

> The outer loop goes row by row and the inner loop goes character by character. This is how all characters end up side by side — we print one row of every character before moving to the next row.
---
