package handlers

import (
	"log"
	"net/http"
)

func Routing() {
	mux := http.NewServeMux()

	// Home page — handles both GET and POST
	mux.HandleFunc("/", HomeHandler)

	// Export endpoint
	mux.HandleFunc("/export", ExportHandler)

	// Serve static files (CSS) under /static/
	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("templates/static"))))

	port := ":8080"
	log.Printf("🌍 Server running at http://localhost%s\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
