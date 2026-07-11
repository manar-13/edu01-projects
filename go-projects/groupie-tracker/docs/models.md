# groupie-tracker — models.go

## Artist

```go
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}
```
Represents one artist or band from the API.
- `ID` — unique number that identifies this artist
- `Image` — URL of the artist's photo
- `Name` — the artist or band name
- `Members` — list of member names
- `CreationDate` — the year the band was formed
- `FirstAlbum` — the date of their first album like `"02-09-1995"`

> The `json:"..."` tags tell Go how to match each field to the JSON key from the API. For example `json:"creationDate"` means the JSON key `creationDate` maps to the `CreationDate` field.

---

## DateEntry and DateData

```go
type DateEntry struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DateData struct {
	Index []DateEntry `json:"index"`
}
```
- `DateEntry` — holds the concert dates for one artist
- `DateData` — wraps a list of all date entries from the API

> The API returns dates inside an `index` array. `DateData` with its `Index` field matches that structure exactly so `json.Unmarshal` can fill it correctly.

---

## RelationEntry and RelationData

```go
type RelationEntry struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationData struct {
	Index []RelationEntry `json:"index"`
}
```
- `RelationEntry` — links concert locations to their dates for one artist
- `DatesLocations` is a map where the key is the location and the value is a list of dates at that location
- `RelationData` — wraps a list of all relation entries from the API

> For example `DatesLocations["paris-france"]` might return `["12-05-2023", "13-05-2023"]`. This is used to show which dates happened at which city on the artist page.

---

## LocationEntry and LocationData

```go
type LocationEntry struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type LocationData struct {
	Index []LocationEntry `json:"index"`
}
```
- `LocationEntry` — holds the list of concert locations for one artist
- `LocationData` — wraps a list of all location entries from the API

> Locations come as raw strings like `"paris-france"` or `"new_york-usa"`. The `cleanLocationName` and `cityLabel` functions in `geo.go` and `main.go` convert these into readable format.

---

## LocOption

```go
type LocOption struct {
	Value string
	Label string
}
```
Represents one option in the location filter dropdown.
- `Value` — the raw API location string like `"paris-france"` — sent to the server when the filter is applied
- `Label` — the clean readable version like `"Paris, France"` — shown to the user in the dropdown

---

## Coordinate

```go
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
```
Holds the geographic coordinates of a location.
- `Lat` — latitude (north-south position)
- `Lng` — longitude (east-west position)

> These two numbers are what Leaflet uses to place a marker on the map. For example Paris is at `Lat: 48.8566, Lng: 2.3522`.

---

## GeocodeResponse

```go
type GeocodeResponse struct {
	Location   string     `json:"location"`
	Coordinate Coordinate `json:"coordinate"`
	Success    bool       `json:"success"`
	Error      string     `json:"error,omitempty"`
}
```
The JSON response our server sends back when the browser asks `/api/geocode`:
- `Location` — the cleaned location name
- `Coordinate` — the lat and lng values
- `Success` — true if geocoding worked, false if it failed
- `Error` — the error message if something went wrong, omitted from JSON if empty

> `omitempty` means if `Error` is an empty string it will not appear in the JSON output at all — keeping the response clean when there is no error.

---

## NominatimResponse

```go
type NominatimResponse []struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}
```
Represents the response from the OpenStreetMap Nominatim API.
- It is a list of results — we always use only the first one
- `Lat` and `Lon` come as strings from Nominatim, not numbers — that is why we parse them with `fmt.Sscanf` in `geo.go`
- `DisplayName` — the full human-readable address returned by Nominatim

---

## SearchSuggestion

```go
type SearchSuggestion struct {
	Label string `json:"label"`
	Type  string `json:"type"`
	ID    int    `json:"id"`
}
```
Represents one suggestion in the search bar dropdown.
- `Label` — the text shown to the user like `"Freddie Mercury - member"`
- `Type` — the category of the match like `"member"`, `"artist/band"`, `"location"`, `"creation date"`, `"first album"`
- `ID` — the artist ID so clicking the suggestion goes to the right artist page
---
