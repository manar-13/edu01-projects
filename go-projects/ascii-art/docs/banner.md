# ascii-art — art/banner.go

## IsASCIIPrintable

```go
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
```
Checks if every character in the text is a normal printable character.
- Loops through every character in the text
- If the character is a backslash `\` skip it — we allow it for `\n`
- If the character number is below 32 or above 126 it is not a standard ASCII character — return false
- If all characters pass the check, return true

> ASCII characters from 32 to 126 cover all normal letters, numbers, spaces, and symbols on a standard keyboard. Anything outside that range (like Arabic or Chinese characters) is rejected.

---

## LoadBanner

```go
func LoadBanner(name string) (map[rune][]string, error) {
	allowed := map[string]bool{
		"standard":   true,
		"shadow":     true,
		"thinkertoy": true,
	}
	if !allowed[name] {
		return nil, fmt.Errorf("invalid banner: %s", name)
	}
```
Takes a banner name and checks if it is one of the three allowed banners. If not, return an error immediately.

---

```go
	file, err := os.Open("banners/" + name + ".txt")
	if err != nil {
		return nil, fmt.Errorf("could not open banner file: %w", err)
	}
	defer file.Close()
```
Opens the banner file. If it cannot be opened, return an error. `defer file.Close()` means the file will be closed automatically when the function finishes.

---

```go
	font := make(map[rune][]string)
	scanner := bufio.NewScanner(file)
	char := rune(32)
	var lines []string
```
Set up everything we need before reading:
- `font` — an empty map that will store each character and its 8 art lines
- `scanner` — a tool to read the file line by line
- `char` — starts at rune 32 which is the space character `' '`
- `lines` — a temporary list to collect the 8 lines of each character

> A rune in Go is a single character. Rune 32 is space, rune 33 is `!`, rune 34 is `"`, and so on. The banner file stores characters in this exact order.

---

```go
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
```
Reads the banner file line by line:
- If the line is empty and we have collected 8 lines — save those 8 lines as the current character in the font map, move to the next character, and reset the lines list
- If the line is not empty — add it to the lines list
- After the loop ends, save the last character if it has 8 lines (it may not be followed by an empty line)
- Return the completed font map

> Each character in the banner file takes exactly 8 lines followed by one empty line. That empty line is the signal that one character ended and the next one is about to begin.

---

## LoadBannerAsLines

```go
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
```
A simpler version of `LoadBanner` — reads the banner file and returns all lines as a plain list without any processing.

> This is used by the color feature which needs the raw lines to work with directly instead of the organized map.
---