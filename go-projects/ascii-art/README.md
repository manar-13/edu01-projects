# ascii-art

## What is ascii-art?

ascii-art is a command-line tool written in Go that converts text into big graphic representations using ASCII characters.

It also supports colors, text alignment, saving output to a file, and reversing ASCII art back into normal text.

---

## Features

| Feature | Description |
|---|---|
| Basic | Convert any text into ASCII art |
| Banner | Choose between 3 different styles |
| Color | Color the whole text or a specific part of it |
| Align | Align the output left, right, center, or justify |
| Output | Save the result into a .txt file |
| Reverse | Convert ASCII art back into normal text |

---

## Banners

There are 3 banner styles available:

| Banner | Description |
|---|---|
| `standard` | Clean and classic style (default) |
| `shadow` | Shadow style |
| `thinkertoy` | Playful and creative style |

---

## Installation

```bash
git clone https://github.com/manar-13/ascii-art.git
cd ascii-art
go mod tidy
```

---

## Usage

### Basic
```bash
go run . "Hello"
go run . "Hello" shadow
go run . "Hello" thinkertoy
```

### Color — whole text
```bash
go run . --color=red "Hello"
```

### Color — specific part
```bash
go run . --color=blue ell "Hello"
```

### Available colors
`black` `red` `green` `yellow` `blue` `magenta` `cyan` `white` `gray`
`orange` `rose` `sky` `lime` `gold` `brightred` `brightgreen` `brightyellow`
`brightblue` `brightmagenta` `brightcyan` `brightwhite`

### Align
```bash
go run . --align=left "Hello World"
go run . --align=right "Hello World"
go run . --align=center "Hello World"
go run . --align=justify "Hello World"
```

### Output — save to file
```bash
go run . --output=result.txt "Hello"
go run . --output=result.txt "Hello" shadow
```

### Reverse — convert ASCII art back to text
```bash
go run . --output=file.txt "Hello"
go run . --reverse=file.txt
```

### Combined
```bash
go run . --align=center "Hello" shadow
go run . --color=green "Hello" thinkertoy
go run . --output=result.txt "Hello" shadow
```

### New lines
```bash
go run . "Hello\nWorld"
go run . "Hello\n\nWorld"
```

---

## File Structure

```
ascii-art/
├── go.mod
├── main.go
├── banners/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
└── art/
    ├── banner.go
    ├── generate.go
    ├── color.go
    ├── align.go
    ├── output.go
    └── reverse.go
```

---

## Author

**Manar Mohamed**
---