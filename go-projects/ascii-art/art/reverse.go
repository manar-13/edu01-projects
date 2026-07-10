package art

import (
	"fmt"
	"os"
	"strings"
)

func ReverseASCII(filename, banner string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	font, err := LoadBanner(banner)
	if err != nil {
		return "", fmt.Errorf("could not load banner: %w", err)
	}

	fileLines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	// Remove trailing empty lines
	for len(fileLines) > 0 && fileLines[len(fileLines)-1] == "" {
		fileLines = fileLines[:len(fileLines)-1]
	}

	if len(fileLines) == 0 {
		return "", nil
	}

	var result strings.Builder
	i := 0

	for i < len(fileLines) {
		// Detect empty line between ascii art blocks
		if fileLines[i] == "" {
			result.WriteString("\n")
			i++
			continue
		}

		// Try to read 8 lines as one ascii art block
		if i+8 > len(fileLines) {
			break
		}

		block := fileLines[i : i+8]
		decoded := decodeBlock(block, font)
		result.WriteString(decoded)
		i += 8

		// Skip the empty line after the block if present
		if i < len(fileLines) && fileLines[i] == "" {
			i++
		}
	}

	return result.String(), nil
}

func decodeBlock(block []string, font map[rune][]string) string {
	if len(block) == 0 {
		return ""
	}

	// Figure out how many characters are in this block
	// by measuring the width of each known character
	var result strings.Builder
	position := 0
	lineLen := len(block[0])

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
			// Check if all 8 rows match at this position
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
			if match {
				result.WriteRune(ch)
				position += charWidth
				matched = true
				break
			}
		}
		if !matched {
			// Skip one character width if nothing matched
			position++
		}
	}

	return result.String()
}
