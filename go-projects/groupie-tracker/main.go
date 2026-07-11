package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var artists []Artist
var relationData RelationData
var locationData LocationData
var dateData DateData

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

func main() {
	var wg sync.WaitGroup
	var fetchErr error
	var mu sync.Mutex

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

	go func() {
		defer wg.Done()
		if err := fetchJSONCtx(ctx, "https://groupietrackers.herokuapp.com/api/relation", &relationData); err != nil {
			mu.Lock()
			fetchErr = fmt.Errorf("failed to fetch relation: %w", err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if err := fetchJSONCtx(ctx, "https://groupietrackers.herokuapp.com/api/locations", &locationData); err != nil {
			mu.Lock()
			fetchErr = fmt.Errorf("failed to fetch locations: %w", err)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if err := fetchJSONCtx(ctx, "https://groupietrackers.herokuapp.com/api/dates", &dateData); err != nil {
			mu.Lock()
			fetchErr = fmt.Errorf("failed to fetch dates: %w", err)
			mu.Unlock()
		}
	}()

	wg.Wait()
	if fetchErr != nil {
		log.Fatal(fetchErr)
	}

	artistByID := make(map[int]Artist, len(artists))
	for _, a := range artists {
		artistByID[a.ID] = a
	}
	relByID := make(map[int]RelationEntry, len(relationData.Index))
	for _, r := range relationData.Index {
		relByID[r.ID] = r
	}
	datesByID := make(map[int]DateEntry, len(dateData.Index))
	for _, d := range dateData.Index {
		datesByID[d.ID] = d
	}
	locsByID := make(map[int]LocationEntry, len(locationData.Index))
	for _, l := range locationData.Index {
		locsByID[l.ID] = l
	}

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
		ay := firstAlbumYear(a.FirstAlbum)
		if ay > 0 {
			if ay < albumMinBound {
				albumMinBound = ay
			}
			if ay > albumMaxBound {
				albumMaxBound = ay
			}
		}
	}
	if creationMaxBound == 0 {
		creationMinBound, creationMaxBound = 1900, 2025
	}
	if albumMaxBound == 0 {
		albumMinBound, albumMaxBound = 1900, 2025
	}

	log.Println("✅ Data loaded (async), server starting...")

	locOptions := buildLocOptions(locationData.Index)

	funcs := template.FuncMap{
		"has":  func(m map[string]bool, k string) bool { return m[k] },
		"qesc": url.QueryEscape,
	}

	tpl := template.Must(template.New("").Funcs(funcs).ParseGlob("templates/*.html"))

	app := &App{
		Tpl:              tpl,
		RelByID:          relByID,
		DatesByID:        datesByID,
		LocsByID:         locsByID,
		ArtistByID:       artistByID,
		LocOptions:       locOptions,
		CreationMinBound: creationMinBound,
		CreationMaxBound: creationMaxBound,
		AlbumMinBound:    albumMinBound,
		AlbumMaxBound:    albumMaxBound,
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("🌐 Listening on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server error:", err)
	}
}

func firstAlbumYear(album string) int {
	if len(album) < 4 {
		return 0
	}
	yearStr := album[len(album)-4:]
	y, _ := strconv.Atoi(yearStr)
	return y
}

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

func formatMembers(members []string) string {
	return strings.Join(members, ", ")
}

func (a *App) renderErrorPage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	data := struct {
		Code    int
		Message string
	}{
		Code:    code,
		Message: message,
	}
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

func (a *App) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	a.renderErrorPage(w, http.StatusNotFound, "Page not found.")
}

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

func formatTitle(s string) string {
	words := strings.Fields(s)
	for i := range words {
		if len(words[i]) > 1 {
			words[i] = strings.ToUpper(words[i][:1]) + strings.ToLower(words[i][1:])
		} else {
			words[i] = strings.ToUpper(words[i])
		}
	}
	return strings.Join(words, " ")
}

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

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	rawQ := r.URL.Query().Get("q")
	rawLoc := r.URL.Query().Get("loc")
	locq := strings.TrimSpace(rawLoc)

	creationMin := atoi0(r.URL.Query().Get("creationMin"))
	creationMax := atoi0(r.URL.Query().Get("creationMax"))
	albumMin := atoi0(r.URL.Query().Get("albumMin"))
	albumMax := atoi0(r.URL.Query().Get("albumMax"))

	if creationMin != 0 && creationMax != 0 && creationMin > creationMax {
		creationMin, creationMax = creationMax, creationMin
	}
	if albumMin != 0 && albumMax != 0 && albumMin > albumMax {
		albumMin, albumMax = albumMax, albumMin
	}

	memberVals := r.URL.Query()["members"]
	locationVals := r.URL.Query()["location"]

	membersSelected := make(map[string]bool, len(memberVals))
	for _, v := range memberVals {
		membersSelected[v] = true
	}

	locationsSelected := make(map[string]bool, len(locationVals))
	for _, v := range locationVals {
		locationsSelected[v] = true
	}

	noFilters := strings.TrimSpace(rawQ) == "" &&
		locq == "" &&
		creationMin == 0 && creationMax == 0 &&
		albumMin == 0 && albumMax == 0 &&
		len(memberVals) == 0 &&
		len(locationVals) == 0

	qry := Query{
		Q:            rawQ,
		CreationMin:  creationMin,
		CreationMax:  creationMax,
		AlbumMin:     albumMin,
		AlbumMax:     albumMax,
		MemberVals:   membersSelected,
		LocationVals: locationVals,
		LocExact:     locq,
	}

	resolveLoc := func(artistID int) []string {
		if le, ok := a.LocsByID[artistID]; ok {
			return le.Locations
		}
		return nil
	}

	pred := BuildPredicate(qry, resolveLoc)

	var filtered []Artist
	if noFilters {
		filtered = artists
	} else {
		filtered = make([]Artist, 0, len(artists))
		for _, ar := range artists {
			if pred(ar) {
				filtered = append(filtered, ar)
			}
		}
	}

	data := struct {
		Artists           []Artist
		Query             string
		LocQuery          string
		AllLocations      []LocOption
		CreationMin       int
		CreationMax       int
		AlbumMin          int
		AlbumMax          int
		MembersSelected   map[string]bool
		LocationsSelected map[string]bool
		RawQuery          string
		CreationMinBound  int
		CreationMaxBound  int
		AlbumMinBound     int
		AlbumMaxBound     int
	}{
		Artists:           filtered,
		Query:             rawQ,
		LocQuery:          locq,
		AllLocations:      a.LocOptions,
		CreationMin:       creationMin,
		CreationMax:       creationMax,
		AlbumMin:          albumMin,
		AlbumMax:          albumMax,
		MembersSelected:   membersSelected,
		LocationsSelected: locationsSelected,
		RawQuery:          r.URL.RawQuery,
		CreationMinBound:  a.CreationMinBound,
		CreationMaxBound:  a.CreationMaxBound,
		AlbumMinBound:     a.AlbumMinBound,
		AlbumMaxBound:     a.AlbumMaxBound,
	}

	a.render(w, "home.html", data)
}

func (a *App) artistHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		a.renderErrorPage(w, http.StatusBadRequest, "Missing artist ID.")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		a.renderErrorPage(w, http.StatusBadRequest, "Invalid artist ID.")
		return
	}

	selectedArtist, ok := a.ArtistByID[id]
	if !ok {
		a.renderErrorPage(w, http.StatusNotFound, "Artist not found.")
		return
	}

	concerts := make(map[string][]string)
	if re, ok := a.RelByID[selectedArtist.ID]; ok {
		concerts = re.DatesLocations
	}

	var artistDates []string
	if de, ok := a.DatesByID[selectedArtist.ID]; ok {
		artistDates = de.Dates
	}

	locations := make([]string, 0)
	for location := range concerts {
		locations = append(locations, location)
	}

	backURL := "/"
	if ret := r.URL.Query().Get("return"); ret != "" {
		backURL = "/?" + ret
	} else if ref := r.Referer(); ref != "" {
		if u, err := url.Parse(ref); err == nil && u.Path == "/" {
			backURL = u.RequestURI()
		}
	}

	data := struct {
		Artist        Artist
		Concerts      map[string][]string
		Dates         []string
		Locations     []string
		MemberDisplay string
		BackURL       string
	}{
		Artist:        selectedArtist,
		Concerts:      concerts,
		Dates:         artistDates,
		Locations:     locations,
		MemberDisplay: formatMembers(selectedArtist.Members),
		BackURL:       backURL,
	}

	a.render(w, "artist.html", data)
}
