package web

import (
	"fmt"
	"strings"
)

func GenerateASCIIArt(input string, font map[rune][]string) (string, error) {
	var result strings.Builder

	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")

	lines := strings.Split(input, "\n")

	for lineNum, line := range lines {
		if lineNum > 0 {
			result.WriteString("\n")
		}

		if len(line) == 0 {
			continue
		}

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
