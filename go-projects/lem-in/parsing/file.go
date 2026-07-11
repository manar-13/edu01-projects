package parsing

import (
	"bufio"
	"fmt"
	"os"

	"lem-in/core"
)

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
