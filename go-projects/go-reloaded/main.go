package main

import (
	"fmt"
	"go-reloaded/functions"
	"os"
	"strings"
)

func ReadLines(filePath string) ([]string, error) {
	content, readErr := os.ReadFile(filePath)

	if readErr != nil {
		return nil, readErr
	}

	textContent := string(content)
	splitLines := strings.Split(textContent, "\n")
	return splitLines, nil
}

func WriteFile(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0644)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please run using: go run . <input.txt> <output.txt>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	if outputPath == "main.go" {
		fmt.Println("Error: You should not use 'main.go' as the output file.")
		os.Exit(1)
	}

	lines, err := ReadLines(inputPath)
	if err != nil {
		fmt.Println("Could not read file:", err)
		os.Exit(1)
	}

	var outputLines []string
	for _, line := range lines {
		line = functions.ConvFormatWithCount(line)
		line = functions.ConvFormatInText(line)
		line = functions.ReplaceHex(line)
		line = functions.ConvBin(line)
		line = functions.FixPunctuSpacing(line)
		line = functions.FixDoublePunctu(line)
		line = functions.FixSingleQuotes(line)
		line = functions.FixAAnGrammar(line)

		outputLines = append(outputLines, line)
	}

	finalText := strings.Join(outputLines, "\n")

	err = WriteFile(outputPath, finalText)
	if err != nil {
		fmt.Println("Could not write result:", err)
		os.Exit(1)
	}

	fmt.Println("Output saved :) ")
}
