# ascii-art — art/align.go

## getTerminalSize

```go
func getTerminalSize() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 80, 24, err
	}
	sizes := strings.Fields(string(out))
	if len(sizes) != 2 {
		return 80, 24, fmt.Errorf("failed to read terminal size")
	}
	width, _ := strconv.Atoi(sizes[1])
	height, _ := strconv.Atoi(sizes[0])
	return width, height, nil
}
```
Asks the terminal how wide and tall it is right now.
- Runs the system command `stty size` which returns something like `"24 80"` (rows and columns)
- If it fails, it returns a default width of 80 and height of 24
- Converts the result from text into numbers and returns them

> This is important for alignment — we need to know the terminal width so we know how much space we have to work with.

---

## getWordPixelLength

```go
func getWordPixelLength(word string, font map[rune][]string) int {
	length := 0
	for _, ch := range word {
		if art, ok := font[ch]; ok {
			length += len(art[0])
		}
	}
	return length
}
```
Calculates how wide a word is in ASCII art pixels.
- Loops through every character in the word
- Looks up that character in the font map
- Adds the width of its first row to the total
- Returns the total width

> Each ASCII art character has a different width. For example `W` is wider than `i`. This function measures the real pixel width, not just the number of letters.

---

## splitWordsToFit

```go
func splitWordsToFit(words []string, font map[rune][]string, width int) [][]string {
	var lines [][]string
	var currentLine []string
	currentWidth := 0
	spaceWidth := len(font[' '][0])

	for _, word := range words {
		wordWidth := getWordPixelLength(word, font)
		additional := wordWidth
		if len(currentLine) > 0 {
			additional += spaceWidth
		}
		if currentWidth+additional > width {
			lines = append(lines, currentLine)
			currentLine = []string{word}
			currentWidth = wordWidth
		} else {
			currentLine = append(currentLine, word)
			currentWidth += additional
		}
	}
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}
	return lines
}
```
Splits the words into groups that fit within the terminal width.
- Measures the width of each word in ASCII art pixels
- If adding the next word would exceed the terminal width, start a new line
- Returns a list of lines where each line is a list of words that fit

> This is like word-wrap in a text editor — but instead of counting letters, we count actual pixel widths of ASCII art characters.

---

## buildASCIIWordsLine

```go
func buildASCIIWordsLine(words []string, font map[rune][]string, spaceWidth int) []string {
	asciiLines := make([]string, 8)
	for wIdx, word := range words {
		for i := 0; i < 8; i++ {
			for _, ch := range word {
				if art, ok := font[ch]; ok {
					asciiLines[i] += art[i]
				}
			}
			if wIdx != len(words)-1 {
				asciiLines[i] += strings.Repeat(" ", spaceWidth)
			}
		}
	}
	return asciiLines
}
```
Builds the 8 rows of ASCII art for a group of words on one line.
- Creates a list of 8 empty strings — one for each row
- For each word, adds the ASCII art of each character row by row
- Adds a space between words but not after the last word
- Returns the 8 completed rows

> Every ASCII art character is 8 rows tall. This function builds all 8 rows at the same time by looping through them together.

---

## buildJustifiedLine

```go
func buildJustifiedLine(words []string, font map[rune][]string, spaceWidth, width int) []string {
	if len(words) == 0 {
		return make([]string, 8)
	}
	baseLines := buildASCIIWordsLine(words, font, spaceWidth)
	baseLen := len(baseLines[0])
	if baseLen >= width || len(words) == 1 {
		return baseLines
	}
	extra := width - baseLen
	gaps := len(words) - 1
	if gaps <= 0 || extra <= 0 {
		return baseLines
	}
	spacePadding := make([]int, gaps)
	for i := 0; i < extra; i++ {
		spacePadding[i%gaps]++
	}
	justified := make([]string, 8)
	for w := 0; w < len(words); w++ {
		for i := 0; i < 8; i++ {
			for _, ch := range words[w] {
				if art, ok := font[ch]; ok {
					justified[i] += art[i]
				}
			}
			if w < len(spacePadding) {
				justified[i] += strings.Repeat(" ", spaceWidth+spacePadding[w])
			}
		}
	}
	return justified
}
```
Spreads words evenly across the full terminal width.
- First builds the line normally and measures its width
- Calculates how many extra spaces are needed to reach the full terminal width
- Distributes those extra spaces evenly between the gaps between words
- Returns the 8 rows with the extra spacing applied

> For example if the terminal is 100 wide, the text is 70 wide, and there are 3 gaps between 4 words — each gap gets about 10 extra spaces so the line fills the full 100.

---

## writeAlignedLines

```go
func writeAlignedLines(builder *strings.Builder, lines []string, padding int) {
	if padding < 0 {
		padding = 0
	}
	for i := 0; i < 8; i++ {
		builder.WriteString(strings.Repeat(" ", padding) + lines[i] + "\n")
	}
}
```
Writes all 8 rows of an ASCII art line with a left padding of spaces.
- If padding is negative, treat it as zero
- For each of the 8 rows, add the padding spaces then the row content then a new line

> This is how left, right, and center alignment work — we just change how many spaces we add before each row.

---

## PrintASCIIAligned

```go
func PrintASCIIAligned(input string, font map[rune][]string, align string) string {
	if len(font) == 0 {
		return "[Error] Empty font map\n"
	}
	width, _, err := getTerminalSize()
	if err != nil || width <= 0 {
		width = 80
	}
	align = strings.ToLower(align)
	words := strings.Fields(input)
	lines := splitWordsToFit(words, font, width)
	spaceWidth := 5
	if spaceArt, ok := font[' ']; ok && len(spaceArt) > 0 {
		spaceWidth = len(spaceArt[0])
	}
	var builder strings.Builder
	for _, lineWords := range lines {
		switch align {
		case "left":
			writeAlignedLines(&builder, buildASCIIWordsLine(lineWords, font, spaceWidth), 0)
		case "center":
			asciiLines := buildASCIIWordsLine(lineWords, font, spaceWidth)
			padding := (width - len(asciiLines[0])) / 2
			writeAlignedLines(&builder, asciiLines, padding)
		case "right":
			asciiLines := buildASCIIWordsLine(lineWords, font, spaceWidth)
			padding := width - len(asciiLines[0])
			writeAlignedLines(&builder, asciiLines, padding)
		case "justify":
			if len(lineWords) == 1 {
				asciiLines := buildASCIIWordsLine(lineWords, font, spaceWidth)
				padding := (width - len(asciiLines[0])) / 2
				writeAlignedLines(&builder, asciiLines, padding)
			} else {
				justified := buildJustifiedLine(lineWords, font, spaceWidth, width)
				writeAlignedLines(&builder, justified, 0)
			}
		default:
			builder.WriteString("[Error] Unknown alignment: " + align + "\n")
		}
	}
	return builder.String()
}
```
The main alignment function — puts everything together.
- If the font is empty, return an error
- Get the terminal width, default to 80 if it fails
- Split the input into words, then into lines that fit the terminal
- Get the space character width from the font
- For each line of words, apply the correct alignment:
  - `left` — no padding, start from the left
  - `center` — add half the remaining space as padding on the left
  - `right` — add all remaining space as padding on the left
  - `justify` — spread words to fill the full width (if only one word, center it instead)
- Returns the complete aligned ASCII art as one string
---