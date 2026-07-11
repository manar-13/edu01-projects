package main

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type DateEntry struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DateData struct {
	Index []DateEntry `json:"index"`
}

type RelationEntry struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationData struct {
	Index []RelationEntry `json:"index"`
}

type LocationEntry struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type LocationData struct {
	Index []LocationEntry `json:"index"`
}

type LocOption struct {
	Value string
	Label string
}

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type GeocodeResponse struct {
	Location   string     `json:"location"`
	Coordinate Coordinate `json:"coordinate"`
	Success    bool       `json:"success"`
	Error      string     `json:"error,omitempty"`
}

type NominatimResponse []struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

type SearchSuggestion struct {
	Label string `json:"label"`
	Type  string `json:"type"`
	ID    int    `json:"id"`
}
