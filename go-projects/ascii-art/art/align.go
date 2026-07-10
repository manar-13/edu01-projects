package art

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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

func getWordPixelLength(word string, font map[rune][]string) int {
	length := 0
	for _, ch := range word {
		if art, ok := font[ch]; ok {
			length += len(art[0])
		}
	}
	return length
}

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

func writeAlignedLines(builder *strings.Builder, lines []string, padding int) {
	if padding < 0 {
		padding = 0
	}
	for i := 0; i < 8; i++ {
		builder.WriteString(strings.Repeat(" ", padding) + lines[i] + "\n")
	}
}

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
