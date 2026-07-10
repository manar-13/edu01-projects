# ascii-art — art/output.go

## WriteToFile

```go
func WriteToFile(content, filename string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
```
Takes the ASCII art text and a file name and writes the text into that file.
- `content` — the ASCII art text we want to save
- `filename` — the name of the file to write into like `result.txt`
- `[]byte(content)` — converts the text into bytes because that is what the file system understands
- `0644` — a permission code meaning the owner can read and write the file, others can only read it
- Returns an error if the writing fails, or `nil` if it succeeds

> This is the only function in this file because saving to a file is a simple one step job. The heavy work of generating the ASCII art is done in `generate.go` before this function is called.
---