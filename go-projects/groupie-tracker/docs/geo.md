# groupie-tracker — geo.go

## Cache Setup

```go
var geocodeCache = make(map[string]Coordinate)
var cacheMutex sync.RWMutex
```
Two global variables set up before any function runs:
- `geocodeCache` — a map that stores already-geocoded locations so we don't ask OpenStreetMap for the same location twice
- `cacheMutex` — a read-write lock that protects the cache when multiple requests try to read or write at the same time

> A mutex is like a lock on a door. `RLock` lets many people read at the same time. `Lock` lets only one person write at a time. This prevents data corruption when multiple map requests come in simultaneously.

---

## geocodeHandler

```go
func geocodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	location := r.URL.Query().Get("location")
	if location == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GeocodeResponse{
			Success: false,
			Error:   "Location parameter is required",
		})
		return
	}
```
Handles requests to `/api/geocode?location=paris-france`:
- Only accepts GET requests — anything else returns 405
- Reads the `location` parameter from the URL
- If no location was given, return a JSON error response

---

```go
	cleanLocation := cleanLocationName(location)

	cacheMutex.RLock()
	if coord, exists := geocodeCache[cleanLocation]; exists {
		cacheMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GeocodeResponse{
			Location:   cleanLocation,
			Coordinate: coord,
			Success:    true,
		})
		return
	}
	cacheMutex.RUnlock()
```
- Cleans the location name first
- Locks the cache for reading and checks if we already have the coordinates for this location
- If found in cache — return the cached result immediately without calling OpenStreetMap
- Unlock the cache after reading

> Checking the cache first is important for performance. The artist page calls `/api/geocode` for every concert location. Without caching, visiting the same artist twice would make double the API calls.

---

```go
	coord, err := geocodeLocation(cleanLocation)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GeocodeResponse{
			Location: cleanLocation,
			Success:  false,
			Error:    err.Error(),
		})
		return
	}

	cacheMutex.Lock()
	geocodeCache[cleanLocation] = coord
	cacheMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GeocodeResponse{
		Location:   cleanLocation,
		Coordinate: coord,
		Success:    true,
	})
}
```
If the location was not in the cache:
- Call `geocodeLocation` to get the coordinates from OpenStreetMap
- If it fails, return a JSON error response
- If it succeeds, lock the cache for writing and save the result
- Return the coordinates as a JSON success response

---

## cleanLocationName

```go
func cleanLocationName(location string) string {
	cleaned := strings.ReplaceAll(location, "_", " ")
	cleaned = strings.ReplaceAll(cleaned, "-", " ")

	parts := strings.Split(cleaned, ",")
	if len(parts) >= 2 {
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(string(parts[i][0])) + strings.ToLower(parts[i][1:])
			}
		}
		cleaned = strings.Join(parts, ", ")
	} else {
		cleaned = strings.TrimSpace(cleaned)
		if len(cleaned) > 0 {
			cleaned = strings.ToUpper(string(cleaned[0])) + strings.ToLower(cleaned[1:])
		}
	}

	return cleaned
}
```
Converts the raw API location string into a clean readable format:
- Replaces underscores and hyphens with spaces
- If the location has multiple parts separated by commas, capitalize the first letter of each part
- If it is a single word, just capitalize the first letter
- Returns the cleaned string

> For example `"paris-france"` becomes `"Paris France"` and `"new_york-usa"` becomes `"New York Usa"`. This makes the location readable for both the map popup and the geocoding search.

---

## geocodeLocation

```go
func geocodeLocation(location string) (Coordinate, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", location)
	params.Add("format", "json")
	params.Add("limit", "1")

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return Coordinate{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "GroupieTracker/1.0")
```
Calls the Nominatim OpenStreetMap API to convert a location name into coordinates:
- Builds the request URL with the location as a query parameter
- `format=json` — get the response as JSON
- `limit=1` — we only need the top result
- Sets the `User-Agent` header — Nominatim requires this to identify who is making the request

---

```go
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Coordinate{}, fmt.Errorf("failed to make geocoding request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coordinate{}, fmt.Errorf("geocoding API returned status: %d", resp.StatusCode)
	}

	var nominatimResp NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&nominatimResp); err != nil {
		return Coordinate{}, fmt.Errorf("failed to decode geocoding response: %v", err)
	}

	if len(nominatimResp) == 0 {
		return Coordinate{}, fmt.Errorf("no results found for location: %s", location)
	}
```
- Sends the request and checks for errors
- Checks the HTTP status — if not 200 OK something went wrong
- Decodes the JSON response into `NominatimResponse`
- If no results were returned the location was not found — return an error

---

```go
	result := nominatimResp[0]
	var lat, lng float64

	if _, err := fmt.Sscanf(result.Lat, "%f", &lat); err != nil {
		return Coordinate{}, fmt.Errorf("failed to parse latitude: %v", err)
	}

	if _, err := fmt.Sscanf(result.Lon, "%f", &lng); err != nil {
		return Coordinate{}, fmt.Errorf("failed to parse longitude: %v", err)
	}

	return Coordinate{Lat: lat, Lng: lng}, nil
}
```
- Takes the first result from the response
- Converts the latitude and longitude from strings to float64 numbers
- Returns them as a `Coordinate` struct

> Nominatim returns coordinates as strings like `"48.8566"` not as numbers. `fmt.Sscanf` parses the string into a float64 so Leaflet can use them to place the marker on the map.
---