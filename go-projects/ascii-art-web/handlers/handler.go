package handlers

import (
	"ascii-art-web/art_web"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type PageData struct {
	Title            string
	Art              string
	Banner           string
	Text             string
	Error            bool
	ErrorMessage     string
	ErrorMessageCode string
}

// HomeHandler serves the homepage at "/"
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

// ArtHandler handles form POST requests to "/"
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

	data := PageData{
		Title:  "Ascii-Art-Web Result",
		Art:    result,
		Banner: banner,
		Text:   text,
		Error:  false,
	}
	RenderTemplate(w, "index.html", data)
}

// ExportHandler handles GET requests to "/export"
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

	// Replace literal \n back to real newlines
	art = strings.ReplaceAll(art, "\\n", "\n")

	content := []byte(art)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"ascii-art.txt\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

// NotFoundHandler renders a custom 404 error page
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	RenderTemplate(w, "errors.html", PageData{
		Title:            "Ascii-Art-Web Error",
		Error:            true,
		ErrorMessage:     "Page Not Found",
		ErrorMessageCode: "404 Not Found",
	})
}

// RenderTemplate loads and executes the template with the given data
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

// RenderError helper sets the appropriate HTTP status and renders error page
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
