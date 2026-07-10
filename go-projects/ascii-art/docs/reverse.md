# ascii-art — art/reverse.go

## ReverseASCII

```go
func ReverseASCII(filename, banner string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}
```
Reads the entire ASCII art file into memory. If the file cannot be read, return an error immediately.

---

```go
	font, err := LoadBanner(banner)
	if err != nil {
		return "", fmt.Errorf("could not load banner: %w", err)
	}
```
Loads the banner font map. We need this to compare the ASCII art in the file against known characters so we can identify them.

> This is the key idea of reverse — we have a dictionary of what every character looks like, and we use it to match what we see in the file.

---

```go
	fileLines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	for len(fileLines) > 0 && fileLines[len(fileLines)-1] == "" {
		fileLines = fileLines[:len(fileLines)-1]
	}

	if len(fileLines) == 0 {
		return "", nil
	}
```
- Convert the file content to text and split it into lines
- Replace `\r\n` (Windows line endings) with `\n` first to avoid issues
- Remove any empty lines at the end of the file
- If the file is completely empty, return an empty string

---

```go
	var result strings.Builder
	i := 0

	for i < len(fileLines) {
		if fileLines[i] == "" {
			result.WriteString("\n")
			i++
			continue
		}
```
Create an empty builder to collect the decoded text. Loop through the file lines. If a line is empty it means there was a `\n` in the original text — add a new line to the result and move on.

---

```go
		if i+8 > len(fileLines) {
			break
		}

		block := fileLines[i : i+8]
		decoded := decodeBlock(block, font)
		result.WriteString(decoded)
		i += 8
```
If there are not enough lines left to form a complete 8-line block, stop. Otherwise take the next 8 lines as one block, decode them back into text, add the result, and move forward by 8 lines.

---

```go
		if i < len(fileLines) && fileLines[i] == "" {
			i++
		}
	}

	return result.String(), nil
}
```
After reading a block, skip the empty line that follows it if there is one. When all blocks are processed, return the decoded text.

---

## decodeBlock

```go
func decodeBlock(block []string, font map[rune][]string) string {
	if len(block) == 0 {
		return ""
	}

	var result strings.Builder
	position := 0
	lineLen := len(block[0])
```
Takes 8 lines of ASCII art and figures out what characters they represent. If the block is empty return an empty string. Start scanning from position 0 and use the length of the first row as the total width to scan.

---

```go
	for position < lineLen {
		matched := false
		for ch, art := range font {
			if len(art) == 0 || len(art[0]) == 0 {
				continue
			}
			charWidth := len(art[0])
			if position+charWidth > lineLen {
				continue
			}
```
Scan across the block from left to right. At each position, try every character in the font map:
- Skip characters with empty art
- Skip characters that would go beyond the end of the line

---

```go
			match := true
			for row := 0; row < 8; row++ {
				if row >= len(block) || position+charWidth > len(block[row]) {
					match = false
					break
				}
				if block[row][position:position+charWidth] != art[row] {
					match = false
					break
				}
			}
```
Check if all 8 rows of this character match the block at the current position. If any row does not match, this is not the right character — move on to the next one.

> This is like placing a stencil over the ASCII art and checking if it fits perfectly. We try every stencil (character) until one fits.

---

```go
			if match {
				result.WriteRune(ch)
				position += charWidth
				matched = true
				break
			}
		}
		if !matched {
			position++
		}
	}

	return result.String()
}
```
If a match is found:
- Add that character to the result
- Move the position forward by the width of that character
- Stop trying other characters and move to the next position

If no match is found at this position, skip one step forward and try again.

When the full width has been scanned, return the decoded text.

> The reverse project is the hardest one because instead of building ASCII art from text, we are doing the opposite — reading ASCII art and figuring out what text it came from. The key tool is the font map which lets us compare what we see against what we know.
---