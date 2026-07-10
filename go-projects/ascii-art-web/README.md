# ascii-art-web

## Description

ascii-art-web is a web application written in Go that brings the ascii-art project to life in the browser. Instead of running commands in the terminal, you can type your text, choose a banner style, and instantly see the ASCII art result on the page. You can also download the result as a text file.

The project combines 4 features in one:

| Feature | Description |
|---|---|
| ascii-art-web | The main web server with a GUI to generate ASCII art |
| stylize | A clean, responsive, and user friendly design |
| dockerize | The app packaged in a Docker container for easy deployment |
| exportfile | A download button to export the ASCII art as a `.txt` file |

---

## Authors

**Manar Mohamed**

---

## Usage

### Run normally with Go

```bash
go run .
```

Then open your browser and go to:
```
http://localhost:8080
```

### Run with Docker

Make sure Docker is installed and running, then:

```bash
bash build.sh
```

Then open your browser and go to:
```
http://localhost:8080
```

---

## How to use the website

1. Type your text in the input field
2. Choose a banner style — Standard, Shadow, or Thinkertoy
3. Click **Generate**
4. See the ASCII art result on the page
5. Click **Download as .txt** to export the result

---

## HTTP Endpoints

| Endpoint | Method | Description |
|---|---|---|
| `/` | GET | Serves the main page |
| `/` | POST | Receives the form and generates ASCII art |
| `/export` | GET | Downloads the ASCII art as a `.txt` file |
| `/static/` | GET | Serves CSS and static files |

---

## HTTP Status Codes

| Code | Meaning |
|---|---|
| 200 | Everything worked |
| 400 | Bad request — missing or invalid input |
| 404 | Page not found |
| 405 | Method not allowed |
| 500 | Internal server error |

---

## Banners

| Banner | Description |
|---|---|
| `standard` | Clean and classic style |
| `shadow` | Shadow style |
| `thinkertoy` | Playful and creative style |

---

## File Structure

```
ascii-art-web/
├── go.mod
├── main.go
├── Dockerfile
├── build.sh
├── banners/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
├── art_web/
│   ├── checker.go
│   └── printer.go
├── handlers/
│   ├── handler.go
│   └── route.go
└── templates/
    ├── index.html
    ├── errors.html
    └── static/
        └── style.css
```

---

## Implementation Details

**How the server works:**
1. Server starts and checks all 3 banner files are valid
2. User opens the browser and sees the main page
3. User fills the form and clicks Generate
4. The form sends a POST request to `/`
5. The server reads the text and banner from the form
6. The server loads the banner font map from the `banners/` folder
7. The server generates the ASCII art line by line, row by row
8. The result is sent back and displayed on the same page
9. The user can click Download to export the result via `/export`

**How ASCII art is generated:**
- Each character in the banner file is 8 lines tall
- The server reads the input text character by character
- For each character it finds the matching 8 lines in the font map
- It prints all characters side by side row by row
- The result is one block of text that looks like big ASCII art letters

**How Docker works:**
- Stage 1 builds the Go binary inside a Go container
- Stage 2 copies only the binary and needed files into a lightweight Alpine container
- The final image is small and clean with no unnecessary files
---