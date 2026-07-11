# groupie-tracker — filter.go

## Query

```go
type Query struct {
	Q            string
	CreationMin  int
	CreationMax  int
	AlbumMin     int
	AlbumMax     int
	MemberVals   map[string]bool
	LocationVals []string
	LocExact     string
}
```
A struct that holds all the filter values the user selected on the home page.
- `Q` — the search text typed in the search bar
- `CreationMin` / `CreationMax` — the creation year range
- `AlbumMin` / `AlbumMax` — the first album year range
- `MemberVals` — a map of selected member counts like `{"3": true, "5": true}`
- `LocationVals` — a list of selected location checkboxes
- `LocExact` — the exact location text typed in the location text input

---

## BuildPredicate

```go
func BuildPredicate(q Query, resolveLoc func(artistID int) []string) func(Artist) bool {
	return func(a Artist) bool {
```
Takes the query and returns a function that checks if a single artist passes all the filters.

> This pattern is called a "predicate" — a function that returns true or false. Instead of filtering inside the handler, we build a reusable checker function and apply it to every artist.

---

```go
		if s := strings.ToLower(strings.TrimSpace(q.Q)); s != "" {
			nameMatch := strings.Contains(strings.ToLower(a.Name), s)

			memberMatch := false
			for _, m := range a.Members {
				if strings.Contains(strings.ToLower(m), s) {
					memberMatch = true
					break
				}
			}

			if !nameMatch && !memberMatch {
				return false
			}
		}
```
Search text check — only runs if the user typed something:
- Converts the search text to lowercase for case-insensitive matching
- Checks if the artist name contains the search text
- Checks if any member name contains the search text
- If neither the name nor any member matches — this artist fails the filter

---

```go
		if q.CreationMin != 0 && a.CreationDate < q.CreationMin {
			return false
		}
		if q.CreationMax != 0 && a.CreationDate > q.CreationMax {
			return false
		}
```
Creation year range check:
- If a minimum was set and the artist's creation date is below it — fail
- If a maximum was set and the artist's creation date is above it — fail
- `!= 0` means only check if the user actually set a value

---

```go
		ay := firstAlbumYear(a.FirstAlbum)

		if q.AlbumMin != 0 && ay < q.AlbumMin {
			return false
		}
		if q.AlbumMax != 0 && ay > q.AlbumMax {
			return false
		}
```
First album year range check — same logic as creation year but for the album date.
- `firstAlbumYear` extracts just the year from the album date string like `"02-09-1995"` → `1995`

---

```go
		if len(q.MemberVals) > 0 {
			if !q.MemberVals[strconv.Itoa(len(a.Members))] {
				return false
			}
		}
```
Member count check — only runs if any checkbox was ticked:
- Converts the artist's actual member count to a string like `3` → `"3"`
- Checks if that string exists in the selected member values map
- If the artist's member count is not in the selected values — fail

---

```go
		if len(q.LocationVals) > 0 || strings.TrimSpace(q.LocExact) != "" {
			artistLocs := resolveLoc(a.ID)

			if len(q.LocationVals) > 0 {
				locMatch := false
				for _, artistLoc := range artistLocs {
					for _, selectedLoc := range q.LocationVals {
						if artistLoc == selectedLoc {
							locMatch = true
							break
						}
					}
					if locMatch {
						break
					}
				}
				if !locMatch {
					return false
				}
			}
```
Location checkbox check — only runs if any location was selected:
- `resolveLoc` is a function passed in that gets the artist's locations by ID
- Loops through the artist's locations and compares each one to the selected locations
- If none of the artist's locations match any selected location — fail

---

```go
			if strings.TrimSpace(q.LocExact) != "" {
				match := false
				for _, loc := range artistLocs {
					if strings.EqualFold(loc, q.LocExact) {
						match = true
						break
					}
				}
				if !match {
					return false
				}
			}
		}

		return true
	}
}
```
Exact location text check — only runs if the user typed something in the exact location input:
- Loops through the artist's locations
- `strings.EqualFold` compares case-insensitively — so `"Paris"` matches `"paris"`
- If none of the artist's locations match the exact text — fail
- If the artist passes all checks — return true

> `return true` at the end means the artist passed every single filter that was active. Only artists that reach this line will be shown on the page.
---
