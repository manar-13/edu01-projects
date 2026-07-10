package art

import (
	"bufio"
	"fmt"
	"os"
)

func IsASCIIPrintable(s string) bool {
	for _, r := range s {
		if r == '\\' {
			continue
		}
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

func LoadBanner(name string) (map[rune][]string, error) {
	allowed := map[string]bool{
		"standard":   true,
		"shadow":     true,
		"thinkertoy": true,
	}
	if !allowed[name] {
		return nil, fmt.Errorf("invalid banner: %s", name)
	}
	file, err := os.Open("banners/" + name + ".txt")
	if err != nil {
		return nil, fmt.Errorf("could not open banner file: %w", err)
	}
	defer file.Close()

	font := make(map[rune][]string)
	scanner := bufio.NewScanner(file)
	char := rune(32)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if len(lines) == 8 {
				font[char] = lines
				char++
				lines = nil
			}
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) == 8 {
		font[char] = lines
	}
	return font, scanner.Err()
}

func LoadBannerAsLines(name string) ([]string, error) {
	file, err := os.Open("banners/" + name + ".txt")
	if err != nil {
		return nil, fmt.Errorf("could not open banner file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
