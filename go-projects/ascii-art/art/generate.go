package art

import (
	"fmt"
	"strings"
)

func GenerateASCIIArt(input string, font map[rune][]string) string {
	var result strings.Builder
	lines := strings.Split(input, "\\n")

	for i, line := range lines {
		if line == "" {
			if i > 0 {
				result.WriteString("\n")
			}
			continue
		}
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

func PrintASCIIArt(input string, font map[rune][]string) {
	fmt.Print(GenerateASCIIArt(input, font))
}
