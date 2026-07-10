# ascii-art — main.go

## Package and Imports

```go
package main

import (
	"ascii-art/art"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)
```
This is the starting point of the program. We import:
- `ascii-art/art` — our own art package where all the functions live
- `fmt` — for printing messages to the terminal
- `os` — for reading the command line arguments
- `path/filepath` — for checking file extensions like `.txt`
- `strings` — for working with text

---

## main

```go
if len(os.Args) < 2 {
    fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
    return
}
```
If the user runs the program with no arguments at all, print the usage message and stop.

---

```go
var (
    colorFlag   string
    outputFile  string
    alignType   string
    reverseFile string
    substr      string
    text        string
    banner      = "standard"
)
```
Create empty boxes to store each possible flag and argument as we read them. Banner defaults to `"standard"` since it is the most common one.

---

```go
args := os.Args[1:]

for i := 0; i < len(args); i++ {
    arg := args[i]
    switch {
```
Skip the first argument (the program name) and loop through the rest. For each argument check what it is.

---

```go
    case strings.HasPrefix(arg, "--color="):
        colorFlag = strings.ToLower(strings.TrimPrefix(arg, "--color="))
```
If the argument starts with `--color=` extract the color name and store it in lowercase. For example `--color=RED` becomes `"red"`.

---

```go
    case strings.HasPrefix(arg, "--output="):
        outputFile = strings.TrimPrefix(arg, "--output=")
        if filepath.Ext(outputFile) != ".txt" {
            fmt.Println("Error: Output file must have .txt extension")
            return
        }
```
If the argument starts with `--output=` extract the file name. If the file does not end with `.txt` print an error and stop.

---

```go
    case strings.HasPrefix(arg, "--align="):
        alignType = strings.TrimPrefix(arg, "--align=")
        if alignType != "left" && alignType != "right" && alignType != "center" && alignType != "justify" {
            fmt.Println("Invalid alignment. Use: left, right, center, justify")
            return
        }
```
If the argument starts with `--align=` extract the alignment type. If it is not one of the four valid options, print an error and stop.

---

```go
    case strings.HasPrefix(arg, "--reverse="):
        reverseFile = strings.TrimPrefix(arg, "--reverse=")
```
If the argument starts with `--reverse=` extract the file name to reverse.

---

```go
    case arg == "standard" || arg == "shadow" || arg == "thinkertoy":
        banner = arg
```
If the argument is one of the three banner names, store it as the chosen banner.

---

```go
    case text == "":
        text = arg
    default:
        if colorFlag != "" && substr == "" {
            substr = text
            text = arg
        }
    }
}
```
If we have not found the text yet, store this argument as the text. If color is set and we get a second non-flag argument, it means the first one was the substring and this one is the actual text — swap them.

> For example `--color=red ell "Hello"` — first `ell` is stored as text, then when `"Hello"` arrives we realize `ell` was the substring and `"Hello"` is the real text.

---

```go
if reverseFile != "" {
    result, err := art.ReverseASCII(reverseFile, banner)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(result)
    return
}
```
If a reverse file was given, run the reverse function and print the result. Stop here — no need to do anything else.

---

```go
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
```
Three validation checks on the text:
- If text is empty, print an error and stop
- If text is just `\n`, print an empty line and stop
- If text contains non-ASCII characters, print an error and stop

---

```go
fontMap, err := art.LoadBanner(banner)
if err != nil {
    fmt.Println("Error:", err)
    return
}
```
Load the chosen banner font into a map. If loading fails, print the error and stop.

---

```go
var output string

switch {
case colorFlag != "":
    output = art.DrawColorASCII(text, colorFlag, substr, fontMap)
case alignType != "":
    output = art.PrintASCIIAligned(text, fontMap, alignType)
default:
    output = art.GenerateASCIIArt(text, fontMap)
}
```
Decide which function to call based on which flag was given:
- Color flag set → draw colored ASCII art
- Align flag set → draw aligned ASCII art
- No special flag → draw normal ASCII art

---

```go
if outputFile != "" {
    if err := art.WriteToFile(output, outputFile); err != nil {
        fmt.Println("Error writing file:", err)
    }
    return
}

fmt.Print(output)
```
If an output file was given, save the result to that file. Otherwise print it directly to the terminal.

> `main.go` is the brain of the whole project. It does not do any ASCII art work itself — it just reads what the user wants, validates it, and sends it to the right function in the `art` package.
---