# groupie-tracker — main.go

## Global Variables

```go
var artists []Artist
var relationData RelationData
var locationData LocationData
var dateData DateData
```
Four global variables that store all the data fetched from the API at startup. They are global so every handler can access them without passing them around.

---

## App Struct

```go
type App struct {
	Tpl        *template.Template
	RelByID    map[int]RelationEntry
	DatesByID  map[int]DateEntry
	LocsByID   map[int]LocationEntry
	ArtistByID map[int]Artist
	LocOptions []LocOption

	CreationMinBound int
	CreationMaxBound int
	AlbumMinBound    int
	AlbumMaxBound    int
}
```
A struct that holds everything the server needs to handle requests:
- `Tpl` — all HTML templates loaded and ready to use
- `RelByID`, `DatesByID`, `LocsByID`, `ArtistByID` — maps that let us look up any artist's data by ID instantly instead of looping through lists every time
- `LocOptions` — the sorted list of all locations for the filter dropdown
- The four bound values — the minimum and maximum years found in the data, used to set the filter slider ranges

---

## Data Fetching

```go
ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
defer cancel()

wg.Add(4)

go func() {
	defer wg.Done()
	if err := fetchJSONCtx(ctx, "https://groupietrackers.herokuapp.com/api/artists", &artists); err != nil {
		mu.Lock()
		fetchErr = fmt.Errorf("failed to fetch artists: %w", err)
		mu.Unlock()
	}
}()
```
Fetches all 4 API endpoints at the same time using goroutines:
- Creates a context with a 12 second timeout — if the API takes longer than that, stop and fail
- `wg.Add(4)` — tells the WaitGroup we are starting 4 goroutines
- Each goroutine fetches one endpoint and stores errors safely using a mutex
- `defer wg.Done()` — marks that goroutine as finished when it returns

> A goroutine is like a lightweight thread. Running all 4 fetches at the same time means the startup is 4x faster than fetching them one by one.

---

```go
wg.Wait()
if fetchErr != nil {
	log.Fatal(fetchErr)
}
```
- `wg.Wait()` — blocks until all 4 goroutines are done
- If any fetch failed, stop the server immediately with the error message

---

## Building Lookup Maps

```go
artistByID := make(map[int]Artist, len(artists))
for _, a := range artists {
	artistByID[a.ID] = a
}
```
Converts the list of artists into a map keyed by ID. Same is done for relations, dates, and locations.

> Looking up an artist by ID in a map is instant. Looking through a list requires checking every item one by one. With 50+ artists this difference is small but it is the correct and professional approach.

---

## Calculating Filter Bounds

```go
creationMinBound, creationMaxBound := 9999, 0
albumMinBound, albumMaxBound := 9999, 0
for _, a := range artists {
	if a.CreationDate > 0 {
		if a.CreationDate < creationMinBound {
			creationMinBound = a.CreationDate
		}
		if a.CreationDate > creationMaxBound {
			creationMaxBound = a.CreationDate
		}
	}
	...
}
```
Loops through all artists to find the earliest and latest creation and album years. These become the min and max values of the filter sliders on the home page.

> Starting `creationMinBound` at 9999 and `creationMaxBound` at 0 is a clever trick. Any real year will be smaller than 9999 (updating the min) and larger than 0 (updating the max).

---

## Template Setup

```go
funcs := template.FuncMap{
	"has":  func(m map[string]bool, k string) bool { return m[k] },
	"qesc": url.QueryEscape,
}

tpl := template.Must(template.New("").Funcs(funcs).ParseGlob("templates/*.html"))
```
Loads all HTML templates and registers two custom functions:
- `has` — checks if a key exists in a map, used in templates to check if a checkbox should be ticked
- `qesc` — URL-encodes a string, used to pass filter state through links

> `template.Must` panics if the templates fail to load. This is intentional — if templates are missing the server should not start at all.

---

## Routing

```go
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		app.homeHandler(w, r)
	case "/artist":
		app.artistHandler(w, r)
	case "/api/geocode":
		geocodeHandler(w, r)
	case "/api/search":
		searchHandler(w, r)
	default:
		app.notFoundHandler(w, r)
	}
})
```
Routes every request to the right handler based on the URL path:
- `/` — home page with artist cards and filters
- `/artist` — individual artist page with map and concerts
- `/api/geocode` — converts a location name to coordinates for the map
- `/api/search` — returns search suggestions as JSON
- anything else — 404 not found page

---

## Server Setup

```go
srv := &http.Server{
	Addr:              ":8080",
	Handler:           mux,
	ReadHeaderTimeout: 5 * time.Second,
	ReadTimeout:       10 * time.Second,
	WriteTimeout:      10 * time.Second,
	IdleTimeout:       60 * time.Second,
}
```
Creates a server with timeouts on every stage of the request:
- `ReadHeaderTimeout` — max time to read the request headers
- `ReadTimeout` — max time to read the full request
- `WriteTimeout` — max time to send the response
- `IdleTimeout` — max time to keep an idle connection open

> These timeouts protect the server from slow or malicious clients that hold connections open forever and waste resources.

---

## Helper Functions

```go
func firstAlbumYear(album string) int {
	if len(album) < 4 {
		return 0
	}
	yearStr := album[len(album)-4:]
	y, _ := strconv.Atoi(yearStr)
	return y
}
```
Extracts just the year from a date string like `"02-09-1995"` by taking the last 4 characters and converting them to an integer.

---

```go
func atoi0(s string) int {
	if s == "" {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}
```
Converts a string to an integer safely. Returns 0 if the string is empty or not a valid number. Used to read filter values from the URL without crashing.

---

```go
func cityLabel(apiLoc string) string {
	parts := strings.Split(apiLoc, "-")
	if len(parts) < 2 {
		city := strings.ReplaceAll(apiLoc, "_", " ")
		return formatTitle(city)
	}
	city := strings.ReplaceAll(parts[0], "_", " ")
	country := strings.ReplaceAll(parts[1], "_", " ")
	return formatTitle(city) + ", " + formatTitle(country)
}
```
Converts the raw API location format like `"new_york-usa"` into a readable label like `"New York, Usa"` for the filter dropdown.

---

```go
func buildLocOptions(index []LocationEntry) []LocOption {
	optMap := make(map[string]LocOption)
	for _, entry := range index {
		for _, loc := range entry.Locations {
			if _, ok := optMap[loc]; !ok {
				optMap[loc] = LocOption{Value: loc, Label: cityLabel(loc)}
			}
		}
	}
	opts := make([]LocOption, 0, len(optMap))
	for _, o := range optMap {
		opts = append(opts, o)
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].Label < opts[j].Label })
	return opts
}
```
Builds the sorted list of unique locations for the filter dropdown:
- Uses a map to avoid duplicate locations
- Converts each raw location to a readable label
- Sorts them alphabetically so the dropdown is easy to read

---

## renderErrorPage and render

```go
func (a *App) renderErrorPage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	...
	if err := a.Tpl.ExecuteTemplate(w, "error.html", data); err != nil {
		_, _ = w.Write([]byte("<h1>Error</h1><p>Something went wrong.</p>"))
	}
}

func (a *App) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.Tpl.ExecuteTemplate(w, name, data); err != nil {
		a.renderErrorPage(w, http.StatusInternalServerError, "Template render error.")
	}
}
```
Two reusable helpers used by every handler:
- `renderErrorPage` — sets the HTTP status code and renders the error template
- `render` — renders any named template with data, falls back to error page if it fails

---

## homeHandler

```go
func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	rawQ := r.URL.Query().Get("q")
	...
	noFilters := strings.TrimSpace(rawQ) == "" && ...

	pred := BuildPredicate(qry, resolveLoc)

	var filtered []Artist
	if noFilters {
		filtered = artists
	} else {
		for _, ar := range artists {
			if pred(ar) {
				filtered = append(filtered, ar)
			}
		}
	}

	a.render(w, "home.html", data)
}
```
Handles the home page:
- Reads all filter values from the URL parameters
- If no filters are active, shows all artists directly without filtering
- If filters are active, builds a predicate and applies it to every artist
- Sends all the data the template needs and renders `home.html`

---

## artistHandler

```go
func (a *App) artistHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	...
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		a.renderErrorPage(w, http.StatusBadRequest, "Invalid artist ID.")
		return
	}

	selectedArtist, ok := a.ArtistByID[id]
	...
	backURL := "/"
	if ret := r.URL.Query().Get("return"); ret != "" {
		backURL = "/?" + ret
	}

	a.render(w, "artist.html", data)
}
```
Handles the individual artist page:
- Reads the `id` from the URL and validates it
- Looks up the artist, concerts, dates, and locations by ID
- Builds the `backURL` — if a `return` parameter exists it means the user came from a filtered page, so the back button goes back to those filters
- Renders `artist.html` with all the data
---