package art

import (
	"fmt"
	"strings"
)

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

func Colorize(segment, color string) string {
	colorCode, ok := ColorMap[strings.ToLower(color)]
	if !ok {
		return segment
	}
	return fmt.Sprintf("%s%s\033[0m", colorCode, segment)
}

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

func isInRange(index int, indices []int, length int) bool {
	for _, start := range indices {
		if index >= start && index < start+length {
			return true
		}
	}
	return false
}

func DrawColorASCII(input, color, substr string, font map[rune][]string) string {
	var result strings.Builder
	lines := strings.Split(input, "\\n")

	for i, line := range lines {
		if line == "" {
			if i > 0 {
				result.WriteString("\n")
			}
			continue
		}
		indices := FindSubStringIndices(line, substr)
		runes := []rune(line)
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
