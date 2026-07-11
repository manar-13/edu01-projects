package main

import (
	"strconv"
	"strings"
)

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

func BuildPredicate(q Query, resolveLoc func(artistID int) []string) func(Artist) bool {
	return func(a Artist) bool {
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

		if q.CreationMin != 0 && a.CreationDate < q.CreationMin {
			return false
		}
		if q.CreationMax != 0 && a.CreationDate > q.CreationMax {
			return false
		}

		ay := firstAlbumYear(a.FirstAlbum)

		if q.AlbumMin != 0 && ay < q.AlbumMin {
			return false
		}
		if q.AlbumMax != 0 && ay > q.AlbumMax {
			return false
		}

		if len(q.MemberVals) > 0 {
			if !q.MemberVals[strconv.Itoa(len(a.Members))] {
				return false
			}
		}

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
