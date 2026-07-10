# ascii-art-web — art_web/checker.go

## IsASCIIPrintable

```go
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
```
Checks if every character in the text is a normal printable character.
- Loops through every character one by one
- If the character is between 32 and 126 it is a normal ASCII character — continue
- If the character is a new line `\n`, tab `\t`, or carriage return `\r` — also allow it
- If anything else is found — return false immediately
- If all characters pass — return true

> This protects the server from receiving text it cannot convert into ASCII art, like Arabic or Chinese characters.

---

## LoadBanner

```go
func LoadBanner(name string) (map[rune][]string, error) {
	path := fmt.Sprintf("banners/%s.txt", name)
	if err := ensureBanner(path, name); err != nil {
		return nil, err
	}
	return parseBannerFile(path)
}
```
Takes a banner name like `"standard"` and loads it into a font map.
- Builds the file path — for example `"banners/standard.txt"`
- Checks the banner file exists and is valid
- Parses the file and returns the font map

---

## ensureBanner

```go
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
```
Checks that the banner file exists and is valid before we try to use it.
- Tries to open the file — if it does not exist, return a clear error
- If it opens successfully, close it and try to parse it
- If parsing fails, the file is corrupted — return an error
- If everything is fine, return nil

> This runs every time the server receives a request so we catch problems early before trying to generate art.

---

## parseBannerFile

```go
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
```
Reads the banner file and builds a map of every character and its 8 art lines.
- Opens the file and reads it line by line
- Starts from rune 32 which is the space character
- When an empty line is found and we have 8 lines collected — save them as one character and move to the next
- If a character does not have exactly 8 lines — return an error
- After reading everything, check that the font has at least 95 characters — if not the file is incomplete
- Return the completed font map

> Each character in the banner file takes exactly 8 lines followed by one empty line. That empty line signals that one character ended and the next begins.

---

## EnsureBanners

```go
func EnsureBanners() error {
	names := []string{"standard", "shadow", "thinkertoy"}
	for _, n := range names {
		if _, err := LoadBanner(n); err != nil {
			return fmt.Errorf("failed banner %q: %w", n, err)
		}
	}
	return nil
}
```
Checks all three banner files when the server starts.
- Loops through all three banner names
- Tries to load each one
- If any one fails — return an error and stop the server

> This runs once at startup in `main.go`. If any banner file is missing or broken the server will not start at all — better to catch it early than fail later during a user request.
---