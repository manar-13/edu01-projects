package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var geocodeCache = make(map[string]Coordinate)
var cacheMutex sync.RWMutex

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
