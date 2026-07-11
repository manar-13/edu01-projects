# groupie-tracker — templates/artist.html

## Head and Imports

```html
<link rel="stylesheet" href="/static/style.css" />
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
```
Loads three things the page needs:
- Our own CSS file for the page styling
- Leaflet CSS — styles for the map
- Leaflet JS — the map library that draws and controls the map

> Leaflet is a free open-source map library. It uses OpenStreetMap tiles to display the map — no API key needed.

---

## Back Link and Artist Info

```html
<a href="{{.BackURL}}" class="back-link">← Back to Home</a>

<div class="artist-container">
    <h1>{{.Artist.Name}}</h1>
    <img src="{{.Artist.Image}}" alt="{{.Artist.Name}}" />
    <p><strong>Members:</strong> {{.MemberDisplay}}</p>
    <p><strong>Creation Date:</strong> {{.Artist.CreationDate}}</p>
    <p><strong>First Album:</strong> {{.Artist.FirstAlbum}}</p>
</div>
```
- The back link takes the user back to wherever they came from — if they had filters applied, it goes back to the filtered results
- `{{.BackURL}}` is filled in by the Go server with the correct URL
- The artist container shows the artist image, members, creation date, and first album
- All values between `{{` and `}}` are filled in by the Go template engine at the time of the request

---

## Map Container

```html
<div id="map-container">
    <div id="map"></div>
    <div id="map-loading">Loading concert locations...</div>
</div>
```
Two divs inside the map container:
- `map` — the empty div where Leaflet will draw the actual map
- `map-loading` — a loading message shown while the locations are being geocoded, hidden when done

---

## Concerts by City

```html
{{if .Concerts}}
<div class="concerts">
    {{range $city, $dates := .Concerts}}
    <div class="concert-city">
        <h3>{{$city}}</h3>
        {{range $dates}}
        <div class="concert-date">{{.}}</div>
        {{end}}
    </div>
    {{end}}
</div>
{{else}}
<p style="text-align:center;">No upcoming concerts available.</p>
{{end}}
```
- If there are concerts, loop through the map of city to dates
- For each city show a card with the city name and all its concert dates
- If there are no concerts show a simple message

> `{{range $city, $dates := .Concerts}}` loops through a map where the key is the city name and the value is a list of dates.

---

## All Concert Dates

```html
{{if .Dates}}
<ul style="text-align:center; padding:0; list-style-type:none;">
    {{range .Dates}}
    <li style="margin:5px 0;">🎵 {{.}}</li>
    {{end}}
</ul>
{{else}}
<p style="text-align:center;">No concert dates listed.</p>
{{end}}
```
Shows all concert dates as a simple list. Each date gets a music note emoji. If there are no dates, show a message.

---

## Map Initialization

```javascript
const map = L.map('map').setView([20, 0], 2);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors'
}).addTo(map);
```
- Creates the map inside the `map` div
- Sets the starting view to coordinates `[20, 0]` — roughly the center of the world
- Zoom level `2` shows the whole world
- Adds the OpenStreetMap tile layer — the actual map images

---

## Location Data

```javascript
const locations = {{.Locations}};
const concerts = {{.Concerts}};
```
The Go template engine fills these variables with the real data from the server before the page is sent to the browser.
- `locations` — a list of location strings like `"paris-france"`
- `concerts` — a map of location to list of dates

---

## Geocoding and Markers

```javascript
const loadingDiv = document.getElementById('map-loading');
let locationsProcessed = 0;
let totalLocations = locations ? locations.length : 0;

if (totalLocations === 0) {
    loadingDiv.textContent = 'No concert locations to display.';
} else {
    const markers = [];

    locations.forEach(async (location) => {
        try {
            const response = await fetch(`/api/geocode?location=${encodeURIComponent(location)}`);
            const data = await response.json();
```
- Count how many locations we need to process
- If there are none, show a message immediately
- Otherwise loop through each location and ask our Go server to convert it to coordinates
- `/api/geocode` is our own endpoint in `geo.go` that calls OpenStreetMap

---

```javascript
            if (data.success) {
                const marker = L.marker([data.coordinate.lat, data.coordinate.lng]).addTo(map);
                markers.push(marker);

                const dates = concerts[location] || [];
                let popupContent = `<div class="map-popup">
                    <h4>${data.location}</h4>
                    <p><strong>Concert Dates:</strong></p>
                    <ul>`;

                dates.forEach(date => {
                    popupContent += `<li>${date}</li>`;
                });

                popupContent += `</ul></div>`;
                marker.bindPopup(popupContent);
```
If geocoding succeeded:
- Place a marker on the map at the returned coordinates
- Build a popup that shows the location name and all its concert dates
- Attach the popup to the marker so it appears when clicked

---

```javascript
        } finally {
            locationsProcessed++;
            if (locationsProcessed === totalLocations) {
                loadingDiv.style.display = 'none';
                if (markers.length > 0) {
                    const group = L.featureGroup(markers);
                    map.fitBounds(group.getBounds().pad(0.1));
                }
            }
        }
    });
}
```
After each location is processed (success or fail):
- Count it as done
- When all locations are processed, hide the loading message
- If any markers were added, zoom the map to fit all of them in view with a small padding

> `fitBounds` automatically adjusts the zoom and position so all markers are visible at once — no matter where in the world the concerts are.
---