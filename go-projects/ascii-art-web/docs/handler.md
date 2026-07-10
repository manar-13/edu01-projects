# ascii-art-web — handlers/handler.go

## PageData

```go
type PageData struct {
	Title            string
	Art              string
	Banner           string
	Text             string
	Error            bool
	ErrorMessage     string
	ErrorMessageCode string
}
```
A struct that holds all the data we send to the HTML templates.
- `Title` — the page title shown in the browser tab
- `Art` — the generated ASCII art result
- `Banner` — the selected banner name
- `Text` — the original input text
- `Error` — true if something went wrong
- `ErrorMessage` — the error description
- `ErrorMessageCode` — the HTTP status code like `"404 Not Found"`

> A struct is like a box that holds multiple related values together. We pass this box to the HTML template so it can display the right content.

---

## HomeHandler

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFoundHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		ArtHandler(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "index.html"),
	)
	if err != nil {
		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title: "ASCII Art Generator",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Render Error: "+err.Error(), http.StatusInternalServerError)
	}
}
```
Handles all requests that come to the `/` route.
- If the path is not exactly `/` — show the 404 page
- If the request is a POST — pass it to `ArtHandler` to generate art
- If the request is a GET — load and show the main page with an empty form
- If the template fails to load or render — return a 500 error

> This one handler covers both GET (showing the page) and POST (generating the art) for the same `/` route.

---

## ArtHandler

```go
func ArtHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		RenderError(w, "Bad Request: Could not parse form data", "400 Bad Request", "", "")
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if text == "" || banner == "" {
		RenderError(w, "Bad Request: Text and banner must be provided", "400 Bad Request", banner, text)
		return
	}
```
Handles the form POST request when the user clicks Generate.
- Parses the form data — if it fails return a 400 error
- Reads the `text` and `banner` values from the form
- If either is empty — return a 400 error

---

```go
	if !web.IsASCIIPrintable(text) {
		RenderError(w, "Bad Request: Only ASCII printable characters allowed", "400 Bad Request", banner, text)
		return
	}

	fontMap, err := web.LoadBanner(banner)
	if err != nil {
		RenderError(w, fmt.Sprintf("Internal Server Error: Failed to load banner: %v", err), "500 Internal Server Error", banner, text)
		return
	}

	result, err := web.GenerateASCIIArt(text, fontMap)
	if err != nil {
		RenderError(w, fmt.Sprintf("Internal Server Error: Failed to generate art: %v", err), "500 Internal Server Error", banner, text)
		return
	}
```
- Check if the text contains only valid ASCII characters — if not return a 400 error
- Load the banner font map — if it fails return a 500 error
- Generate the ASCII art — if it fails return a 500 error

---

```go
	data := PageData{
		Title:  "Ascii-Art-Web Result",
		Art:    result,
		Banner: banner,
		Text:   text,
		Error:  false,
	}
	RenderTemplate(w, "index.html", data)
}
```
If everything worked, build a `PageData` with the result and render the `index.html` template with the ASCII art displayed on the page.

---

## ExportHandler

```go
func ExportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RenderError(w, "Method Not Allowed", "405 Method Not Allowed", "", "")
		return
	}

	art := r.URL.Query().Get("art")
	if art == "" {
		RenderError(w, "Bad Request: No art to export", "400 Bad Request", "", "")
		return
	}

	art = strings.ReplaceAll(art, "\\n", "\n")

	content := []byte(art)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"ascii-art.txt\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
```
Handles the download request when the user clicks the export button.
- Only accepts GET requests — if anything else return a 405 error
- Reads the `art` value from the URL query — if empty return a 400 error
- Fixes any escaped new lines back to real new lines
- Sets the HTTP headers to tell the browser this is a file download:
  - `Content-Type` — the file is plain text
  - `Content-Disposition` — tells the browser to download it as `ascii-art.txt`
  - `Content-Length` — the size of the file in bytes
- Writes the file content to the response

> `Content-Disposition: attachment` is what triggers the browser to download the file instead of displaying it on the page.

---

## NotFoundHandler

```go
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	RenderTemplate(w, "errors.html", PageData{
		Title:            "Ascii-Art-Web Error",
		Error:            true,
		ErrorMessage:     "Page Not Found",
		ErrorMessageCode: "404 Not Found",
	})
}
```
Sends a 404 status code and renders the error page when a route is not found.

---

## RenderTemplate

```go
func RenderTemplate(w http.ResponseWriter, tmplName string, data PageData) {
	tmplPath := filepath.Join("templates", tmplName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Internal Server Error: Failed to parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("TEMPLATE EXECUTION ERROR:", err)
		http.Error(w, "Internal Server Error: Failed to render template", http.StatusInternalServerError)
	}
}
```
A reusable helper that loads an HTML template file and fills it with data.
- Builds the full path to the template file inside the `templates/` folder
- Parses the template file — if it fails return a 500 error
- Executes the template with the given data — if it fails return a 500 error

> Every handler uses this function instead of repeating the same template loading code. This keeps the code clean and consistent.

---

## RenderError

```go
func RenderError(w http.ResponseWriter, errorMsg string, errorCode string, banner string, text string) {
	switch errorCode {
	case "400 Bad Request":
		w.WriteHeader(http.StatusBadRequest)
	case "404 Not Found":
		w.WriteHeader(http.StatusNotFound)
	case "405 Method Not Allowed":
		w.WriteHeader(http.StatusMethodNotAllowed)
	case "500 Internal Server Error":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	data := PageData{
		Title:            "Ascii-Art-Web Error",
		ErrorMessage:     errorMsg,
		ErrorMessageCode: errorCode,
		Banner:           banner,
		Text:             text,
		Error:            true,
	}
	RenderTemplate(w, "errors.html", data)
}
```
A reusable helper that sets the correct HTTP status code and renders the error page.
- Checks the error code and sets the matching HTTP status
- Builds a `PageData` with the error details
- Renders the `errors.html` template with that data

> Every error in every
---