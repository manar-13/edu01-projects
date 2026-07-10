package main

import (
	"ascii-art/art"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
		return
	}

	var (
		colorFlag   string
		outputFile  string
		alignType   string
		reverseFile string
		substr      string
		text        string
		banner      = "standard"
	)

	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case strings.HasPrefix(arg, "--color="):
			colorFlag = strings.ToLower(strings.TrimPrefix(arg, "--color="))
		case strings.HasPrefix(arg, "--output="):
			outputFile = strings.TrimPrefix(arg, "--output=")
			if filepath.Ext(outputFile) != ".txt" {
				fmt.Println("Error: Output file must have .txt extension")
				return
			}
		case strings.HasPrefix(arg, "--align="):
			alignType = strings.TrimPrefix(arg, "--align=")
			if alignType != "left" && alignType != "right" && alignType != "center" && alignType != "justify" {
				fmt.Println("Invalid alignment. Use: left, right, center, justify")
				return
			}
		case strings.HasPrefix(arg, "--reverse="):
			reverseFile = strings.TrimPrefix(arg, "--reverse=")
		case arg == "standard" || arg == "shadow" || arg == "thinkertoy":
			banner = arg
		case text == "":
			text = arg
		default:
			if colorFlag != "" && substr == "" {
				substr = text
				text = arg
			}
		}
	}

	// Handle reverse
	if reverseFile != "" {
		result, err := art.ReverseASCII(reverseFile, banner)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(result)
		return
	}

	// Validate text
	if text == "" {
		fmt.Println("Error: Missing input text")
		return
	}

	if text == `\n` {
		fmt.Println()
		return
	}

	if !art.IsASCIIPrintable(text) {
		fmt.Println("Error: English letters only")
		return
	}

	// Load banner
	fontMap, err := art.LoadBanner(banner)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var output string

	switch {
	case colorFlag != "":
		output = art.DrawColorASCII(text, colorFlag, substr, fontMap)
	case alignType != "":
		output = art.PrintASCIIAligned(text, fontMap, alignType)
	default:
		output = art.GenerateASCIIArt(text, fontMap)
	}

	if outputFile != "" {
		if err := art.WriteToFile(output, outputFile); err != nil {
			fmt.Println("Error writing file:", err)
		}
		return
	}

	fmt.Print(output)
}
