# lem-in — parsing/file.go

## ReadFile

```go
func ReadFile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		core.Fatal(fmt.Sprintf("failed to open file: %s", path))
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		core.Fatal(fmt.Sprintf("error while reading file: %s", err.Error()))
	}

	if len(lines) == 0 {
		core.Fatal("input file is empty")
	}

	return lines
}
```
Opens the input file and reads it line by line into a list of strings.
- Opens the file — if it fails, print an error and stop
- `defer f.Close()` — closes the file automatically when the function finishes
- Reads each line using a buffered scanner
- If the scanner encounters an error, stop with an error message
- If the file is empty, stop with an error message
- Returns all lines as a list of strings

> The raw lines are stored in `data.Input` and printed back to the terminal before the moves are shown — this is how the program echoes the input file as required by the project.
---
