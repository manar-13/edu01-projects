# go-reloaded — main.go

## Package and Imports

```go
package main
```
This file is the starting point of the program. Every Go program must have one `package main`.

---

```go
import (
    "fmt"
    "go-reloaded/functions"
    "os"
    "strings"
)
```
We are borrowing tools from other packages:
- `fmt` — for printing messages to the terminal
- `go-reloaded/functions` — our own functions file we created
- `os` — for reading and writing files
- `strings` — for working with text

---

## ReadLines

```go
func ReadLines(filePath string) ([]string, error) {
```
A function that takes a file path and returns two things: a list of lines from the file, and an error if something goes wrong.

> In Go, functions can return two things at once. The second one is usually an error.

---

```go
content, readErr := os.ReadFile(filePath)

if readErr != nil {
    return nil, readErr
}
```
Read the file and store everything in `content`. If reading fails, return the error immediately and stop.

> `nil` means "nothing". So `return nil, readErr` means "return no lines, but return the error."

---

```go
textContent := string(content)
splitLines := strings.Split(textContent, "\n")
return splitLines, nil
```
- Convert the file content into text
- Split it into separate lines wherever there is a new line `\n`
- Return the lines and `nil` for error — meaning no error happened

---

## WriteFile

```go
func WriteFile(path string, data string) error {
    return os.WriteFile(path, []byte(data), 0644)
}
```
Takes a file path and text, then writes that text into the file.
- `[]byte(data)` — converts text into bytes because that is what the file system understands
- `0644` — a permission code meaning the owner can read and write this file

---

## main

```go
if len(os.Args) != 3 {
    fmt.Println("Please run using: go run . <input.txt> <output.txt>")
    os.Exit(1)
}
```
When you run the program you type `go run . sample.txt result.txt`. `os.Args` stores everything you typed as a list. If the list does not have exactly 3 items, print an error and stop.

> `os.Exit(1)` means "stop the program right now, something went wrong."

---

```go
inputPath := os.Args[1]
outputPath := os.Args[2]
```
Store the two file names in boxes for easy use later.

---

```go
if outputPath == "main.go" {
    fmt.Println("Error: You should not use 'main.go' as the output file.")
    os.Exit(1)
}
```
A safety check — if someone accidentally tries to write the result into `main.go`, stop them before they destroy the code file.

---

```go
lines, err := ReadLines(inputPath)
if err != nil {
    fmt.Println("Could not read file:", err)
    os.Exit(1)
}
```
Call `ReadLines` to read the input file. If it fails, print the error and stop.

---

```go
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
```
- Create an empty list called `outputLines`
- Loop through every line from the input file
- Pass each line through all 8 functions one by one — each function fixes something
- Add the fixed line to `outputLines`

> The order matters. `ConvFormatWithCount` runs before `FixPunctuSpacing` because we handle the instructions before fixing punctuation.

---

```go
finalText := strings.Join(outputLines, "\n")
```
Join all the fixed lines back together into one big text, with a new line between each one.

---

```go
err = WriteFile(outputPath, finalText)
if err != nil {
    fmt.Println("Could not write result:", err)
    os.Exit(1)
}

fmt.Println("Output saved :) ")
```
Write the final text into the output file. If it fails, print the error and stop. If it works, print `Output saved :)`.
---