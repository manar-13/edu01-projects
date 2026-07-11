package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchSuggestion{})
		return
	}

	var suggestions []SearchSuggestion
	seen := make(map[string]bool)

	for _, a := range artists {
		// Match artist/band name
		if strings.Contains(strings.ToLower(a.Name), query) {
			key := a.Name + "-artist/band"
			if !seen[key] {
				suggestions = append(suggestions, SearchSuggestion{
					Label: a.Name + " - artist/band",
					Type:  "artist/band",
					ID:    a.ID,
				})
				seen[key] = true
			}
		}

		// Match members
		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), query) {
				key := m + "-member"
				if !seen[key] {
					suggestions = append(suggestions, SearchSuggestion{
						Label: m + " - member",
						Type:  "member",
						ID:    a.ID,
					})
					seen[key] = true
				}
			}
		}

		// Match first album date
		if strings.Contains(strings.ToLower(a.FirstAlbum), query) {
			key := a.FirstAlbum + "-first album"
			if !seen[key] {
				suggestions = append(suggestions, SearchSuggestion{
					Label: a.FirstAlbum + " - first album",
					Type:  "first album",
					ID:    a.ID,
				})
				seen[key] = true
			}
		}

		// Match creation date
		creationStr := strings.TrimSpace(strings.Join(strings.Fields(
			strings.ReplaceAll(string(rune('0'+a.CreationDate/1000)), "", "")),
			""))
		creationStr = formatInt(a.CreationDate)
		if strings.Contains(creationStr, query) {
			key := creationStr + "-creation date"
			if !seen[key] {
				suggestions = append(suggestions, SearchSuggestion{
					Label: creationStr + " - creation date",
					Type:  "creation date",
					ID:    a.ID,
				})
				seen[key] = true
			}
		}
	}

	// Match locations
	for _, entry := range locationData.Index {
		for _, loc := range entry.Locations {
			cleanLoc := cleanLocationName(loc)
			if strings.Contains(strings.ToLower(cleanLoc), query) {
				key := cleanLoc + "-location"
				if !seen[key] {
					suggestions = append(suggestions, SearchSuggestion{
						Label: cleanLoc + " - location",
						Type:  "location",
						ID:    entry.ID,
					})
					seen[key] = true
				}
			}
		}
	}

	// Limit to 10 suggestions
	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

func formatInt(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
