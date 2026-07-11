# groupie-tracker

## Description

Groupie Tracker is a web application written in Go that fetches data from a music API and displays information about artists and bands in a clean, interactive website. You can search, filter, and explore artists and their concert locations on an interactive map.

The project combines 5 features in one:

| Feature | Description |
|---|---|
| groupie-tracker | The main web server that fetches and displays artist data |
| filters | Filter artists by creation date, album date, member count, and location |
| geolocalization | Shows concert locations on an interactive map using OpenStreetMap |
| visualizations | A clean, responsive, and consistent dark mode design |
| search-bar | A live search bar with typed suggestions by name, member, location, and date |

---

## Authors

**Manar Mohamed**

---

## Usage

### Run the server

```bash
go run .
```

Then open your browser and go to:
```
http://localhost:8080
```

---

## How to use the website

1. Open the home page to see all artists
2. Use the search bar to find an artist, member, location, or date — suggestions appear as you type
3. Use the filters to narrow down the results:
   - Slide the creation year range to filter by when a band was formed
   - Slide the first album year range to filter by album date
   - Tick member count checkboxes to filter by band size
   - Open the locations dropdown to filter by concert location
4. Click Apply Filters to see the results
5. Click an artist card to open their detail page
6. The artist page shows their info, an interactive concert map, concerts by city, and all concert dates
7. Click the back button to return to your filtered results

---

## HTTP Endpoints

| Endpoint | Method | Description |
|---|---|---|
| `/` | GET | Home page with artist cards, search bar and filters |
| `/artist` | GET | Artist detail page with map and concert info |
| `/api/geocode` | GET | Converts a location name to coordinates for the map |
| `/api/search` | GET | Returns search suggestions as JSON |
| `/static/` | GET | Serves CSS and static files |

---

## API Used

Data is fetched from:

| Endpoint | What it contains |
|---|---|
| `/api/artists` | Artist names, images, members, creation date, first album |
| `/api/locations` | Concert locations for each artist |
| `/api/dates` | Concert dates for each artist |
| `/api/relation` | Links locations to their concert dates per artist |

All 4 endpoints are fetched at the same time when the server starts using goroutines for fast loading.

---

## Features Details

### Search Bar
- Live suggestions appear as you type
- Matches by artist name, member name, location, first album date, and creation date
- Each suggestion shows its type — for example `Freddie Mercury - member`
- Click a suggestion to go directly to that artist page
- Use arrow keys to navigate suggestions and Enter to select

### Filters
- Creation year — range slider and number input for min and max
- First album year — range slider and number input for min and max
- Number of members — checkboxes from 1 to 6
- Locations — collapsible dropdown with all unique concert locations plus an exact text search
- Filters stay active when you go to an artist page and come back

### Map
- Each artist page shows an interactive OpenStreetMap with markers for every concert location
- Click a marker to see the location name and all concert dates at that location
- The map automatically zooms to fit all markers in view

---

## File Structure

```
groupie-tracker/
├── go.mod
├── main.go
├── models.go
├── fetch.go
├── filter.go
├── geo.go
├── search.go
├── static/
│   └── style.css
└── templates/
    ├── home.html
    ├── artist.html
    └── error.html
```

---

## Implementation Details

**How the server works:**
1. Server starts and fetches all 4 API endpoints at the same time using goroutines
2. Data is stored in lookup maps keyed by artist ID for fast access
3. Filter bounds (min/max years) are calculated from the real data
4. All HTML templates are loaded once at startup
5. Every request is routed to the right handler based on the URL path

**How filters work:**
- All filter values come from URL parameters
- A predicate function is built from the active filters
- Every artist is checked against the predicate
- Only artists that pass all active filters are shown
- Filter state is preserved in the URL so the back button works correctly

**How the map works:**
- The artist page sends each concert location to `/api/geocode`
- The geocode endpoint calls OpenStreetMap Nominatim to get coordinates
- Results are cached in memory so the same location is never geocoded twice
- Leaflet.js places markers on the map using the returned coordinates

**How search works:**
- As the user types, the browser calls `/api/search?q=...`
- The server checks every artist name, member, album date, creation date, and location
- Matching results are returned as JSON with a label and type
- The browser builds the dropdown from the JSON response
---
