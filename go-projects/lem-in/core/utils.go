package core

import (
	"fmt"
	"os"
	"strconv"
)

func Fatal(msg string) {
	if msg == "" {
		msg = "invalid data format"
	}
	fmt.Printf("ERROR: %s\n", msg)
	os.Exit(1)
}

func Atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		Fatal(fmt.Sprintf("invalid number: '%s'", s))
	}
	return n
}

func IsNumber(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.Atoi(s)
	return err == nil
}

func IsComment(line string) bool {
	return len(line) >= 1 && line[0] == '#' && (len(line) == 1 || line[1] != '#')
}

func IsUnknownCommand(line string) bool {
	if len(line) < 2 {
		return false
	}
	if line[0] != '#' || line[1] != '#' {
		return false
	}
	if line == "##start" || line == "##end" {
		return false
	}
	return true
}

func IsValidRoomName(name string) bool {
	if name == "" {
		return false
	}

	for _, r := range name {
		if r < 32 || r > 126 {
			return false
		}
	}

	return true
}
