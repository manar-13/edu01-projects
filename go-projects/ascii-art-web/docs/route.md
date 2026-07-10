# ascii-art-web — handlers/route.go

## Routing

```go
func Routing() {
	mux := http.NewServeMux()
```
Creates a new router called `mux`. Think of it like a traffic controller — it looks at every request that comes in and decides which handler should deal with it.

---

```go
	mux.HandleFunc("/", HomeHandler)
```
Any request to `/` goes to `HomeHandler`.
- GET `/` — shows the main page
- POST `/` — generates the ASCII art

---

```go
	mux.HandleFunc("/export", ExportHandler)
```
Any request to `/export` goes to `ExportHandler` which sends the ASCII art as a downloadable file.

---

```go
	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("templates/static"))))
```
Serves all static files like CSS from the `templates/static/` folder.
- Any request that starts with `/static/` is handled here
- `http.StripPrefix` removes the `/static/` part from the path before looking for the file
- `http.FileServer` serves the actual file from the `templates/static` folder

> For example a request to `/static/style.css` becomes a lookup for `style.css` inside `templates/static/`. This is how the browser loads the CSS.

---

```go
	port := ":8080"
	log.Printf("🌍 Server running at http://localhost%s\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
```
- Sets the port to `8080`
- Prints a message so you know the server is running
- Starts the server and keeps it running — if it fails for any reason, print the error and stop

> `http.ListenAndServe` blocks forever — it keeps listening for requests until the program is stopped or crashes. That is why it is always the last line in the routing function.
---