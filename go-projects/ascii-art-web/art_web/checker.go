package web

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func IsASCIIPrintable(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 32 && c <= 126 {
			continue
		}
		if c == '\n' || c == '\t' || c == '\r' {
			continue
		}
		return false
	}
	return true
}

func LoadBanner(name string) (map[rune][]string, error) {
	path := fmt.Sprintf("banners/%s.txt", name)
	if err := ensureBanner(path, name); err != nil {
		return nil, err
	}
	return parseBannerFile(path)
}

func ensureBanner(path, name string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("banner file not found: %s", path)
		}
		return err
	}
	f.Close()

	if _, err := parseBannerFile(path); err != nil {
		return fmt.Errorf("banner %q is invalid: %v", name, err)
	}
	return nil
}

func parseBannerFile(path string) (map[rune][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	font := make(map[rune][]string)
	currentRune := rune(32)
	var lines []string

	for scanner.Scan() {
		txt := scanner.Text()
		if txt == "" {
			if len(lines) > 0 {
				if len(lines) != 8 {
					return nil, fmt.Errorf("banner char %q has %d lines", currentRune, len(lines))
				}
				font[currentRune] = append([]string{}, lines...)
				currentRune++
				lines = nil
			}
			continue
		}
		lines = append(lines, txt)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) > 0 {
		if len(lines) != 8 {
			return nil, fmt.Errorf("banner char %q has %d lines", currentRune, len(lines))
		}
		font[currentRune] = append([]string{}, lines...)
	}

	if len(font) < 95 {
		return nil, errors.New("banner missing ASCII characters")
	}
	return font, nil
}

func EnsureBanners() error {
	names := []string{"standard", "shadow", "thinkertoy"}
	for _, n := range names {
		if _, err := LoadBanner(n); err != nil {
			return fmt.Errorf("failed banner %q: %w", n, err)
		}
	}
	return nil
}
