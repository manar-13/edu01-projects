# groupie-tracker — search.go

## searchHandler

```go
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
```
Handles requests to `/api/search?q=...` from the search bar:
- Only accepts GET requests — anything else returns 405
- Reads the `q` parameter and converts it to lowercase for case-insensitive matching
- Trims any extra spaces from the query
- If the query is empty, return an empty list immediately

---

```go
	var suggestions []SearchSuggestion
	seen := make(map[string]bool)
```
Two variables to build the results:
- `suggestions` — the list of matches we will return
- `seen` — a map to track what we already added so we never add the same suggestion twice

---

```go
	for _, a := range artists {
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
```
Loops through every artist and checks if the artist name contains the query:
- Creates a unique key combining the name and type
- If we have not seen this key before, add it to suggestions and mark it as seen
- The label shown to the user is `"Queen - artist/band"`

---

```go
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
```
Checks every member of the current artist against the query:
- Loops through the members list
- If a member name contains the query, add a suggestion
- The label shown to the user is `"Freddie Mercury - member"`
- The ID points to the artist this member belongs to so clicking goes to the right page

---

```go
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
```
Checks if the first album date contains the query:
- Useful when a user types a year like `"1995"` to find artists by album date
- The label shown to the user is `"02-09-1995 - first album"`

---

```go
		creationStr := formatInt(a.CreationDate)
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
```
Converts the creation year to a string and checks if it contains the query:
- `formatInt` converts the integer year to a string like `1970`
- Useful when a user types a year to find artists by when they were formed
- The label shown to the user is `"1970 - creation date"`

---

```go
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
```
Checks all concert locations from the location data:
- Cleans the raw API location string first using `cleanLocationName`
- Checks if the cleaned location contains the query
- The label shown to the user is `"Paris France - location"`
- The ID points to the artist who has a concert at that location

---

```go
	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}
```
- Limits the results to 10 suggestions maximum so the dropdown does not get too long
- Sets the response type to JSON
- Encodes the suggestions list as JSON and sends it to the browser

> The browser's `fetchSuggestions` function in `home.html` receives this JSON and builds the dropdown list from it.

---

## formatInt

```go
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
```
Converts an integer to a string without using `strconv` or `fmt`:
- If the number is 0 return `"0"` immediately
- Extracts the last digit of the number using `n % 10`
- Converts the digit to its character equivalent using `rune('0' + digit)`
- Adds it to the front of the result string
- Removes the last digit from the number using `n / 10`
- Repeats until no digits remain

> For example `1970` → last digit is `0` → result is `"0"` → then `197` → last digit is `7` → result is `"70"` → and so on until result is `"1970"`.
---
