# groupie-tracker — templates/home.html

## Head and Imports

```html
<link rel="stylesheet" href="/static/style.css" />
```
Loads our CSS file for all the page styling. No external libraries needed on the home page.

---

## Search Bar

```html
<div style="position:relative; display:inline-block; width:60%; max-width:520px;">
    <input id="q" type="text" name="q"
        placeholder="Search artist, member, location..."
        value="{{.Query}}"
        autocomplete="off"
        oninput="fetchSuggestions(this.value)"
        onkeydown="handleKey(event)">
    <ul id="suggestions" style="display:none; position:absolute; ..."></ul>
</div>
```
- The input field is wrapped in a relative div so the suggestions dropdown can be positioned directly below it
- `value="{{.Query}}"` — keeps the search text after the form is submitted
- `autocomplete="off"` — disables the browser's own autocomplete so our custom suggestions show instead
- `oninput` — calls `fetchSuggestions` every time the user types a character
- `onkeydown` — calls `handleKey` to handle arrow keys, enter, and escape
- The `ul` is the suggestions dropdown — hidden by default, shown by JavaScript

---

## Creation Year Filter

```html
<div class="range-box">
    <strong>Creation Year</strong>
    <input type="number" name="creationMin"
        min="{{.CreationMinBound}}" max="{{.CreationMaxBound}}"
        value="{{if .CreationMin}}{{.CreationMin}}{{else}}{{.CreationMinBound}}{{end}}">
    <input type="range" name="creationMin" ...>

    <input type="number" name="creationMax" ...>
    <input type="range" name="creationMax" ...>
    <div class="hint">Range: {{.CreationMinBound}}–{{.CreationMaxBound}}</div>
</div>
```
Two pairs of inputs — one number box and one slider for min and max:
- Both inputs share the same `name` so the form sends one value per bound
- The number and range inputs mirror each other visually
- `{{if .CreationMin}}` — if the user already set a value keep it, otherwise use the default bound
- The hint shows the full available range to the user

> The same pattern is repeated for First Album Year with `albumMin` and `albumMax`.

---

## Members Filter

```html
<div class="range-box">
    <strong>Number of Members</strong>
    <label><input type="checkbox" name="members" value="1"
        {{if has .MembersSelected "1"}}checked{{end}}> 1</label>
    <label><input type="checkbox" name="members" value="2"
        {{if has .MembersSelected "2"}}checked{{end}}> 2</label>
    ...
</div>
```
Six checkboxes — one for each possible number of members from 1 to 6.
- `{{if has .MembersSelected "1"}}checked{{end}}` — keeps the checkbox ticked after the form is submitted
- `has` is a custom Go template function that checks if a key exists in a map

---

## Locations Filter

```html
<details class="loc-details" {{if .LocationsSelected}}open{{end}}>
    <summary class="loc-summary">
        <span>Locations</span>
        <span class="arrow"></span>
    </summary>
    <input type="text" name="loc" value="{{.LocQuery}}"
        placeholder="e.g., Seattle, Washington, USA">
    <div class="loc-list">
        {{range .AllLocations}}
        <label class="loc-item">
            <input type="checkbox" name="location" value="{{.Value}}"
                {{if has $.LocationsSelected .Value}}checked{{end}}>
            {{.Label}}
        </label>
        {{end}}
    </div>
</details>
```
A collapsible dropdown using the HTML `details` and `summary` elements — no JS needed:
- Opens automatically if any location was already selected
- Has a text input for an exact location search
- Lists all available locations as checkboxes
- `{{range .AllLocations}}` loops through all unique locations from the API
- The arrow rotates when the dropdown opens using CSS only

---

## Filter Buttons

```html
<button type="submit">Apply Filters</button>
<a href="/">Reset</a>
```
- Submit sends all filter values to the server as URL parameters
- Reset is just a link back to `/` which clears all filters

---

## Artist Cards

```html
{{range .Artists}}
<li>
    <a href="/artist?id={{.ID}}&return={{qesc $.RawQuery}}">
        <img src="{{.Image}}" alt="{{.Name}}">
        <p>{{.Name}}</p>
    </a>
</li>
{{else}}
<li class="no-results">No artists found.</li>
{{end}}
```
- Loops through the filtered list of artists and shows each one as a card
- The link to the artist page includes `return={{qesc $.RawQuery}}` — this saves the current filters so the back button on the artist page returns to the same filtered results
- `qesc` is a custom Go template function that URL-encodes the query string
- `{{else}}` runs when the list is empty — shows a "No artists found" message

---

## fetchSuggestions

```javascript
async function fetchSuggestions(query) {
    const box = document.getElementById('suggestions');
    if (!query || query.length < 1) {
        box.style.display = 'none';
        box.innerHTML = '';
        return;
    }

    try {
        const res = await fetch('/api/search?q=' + encodeURIComponent(query));
        const data = await res.json();
        box.innerHTML = '';
        selectedIndex = -1;

        if (!data || data.length === 0) {
            box.style.display = 'none';
            return;
        }

        data.forEach((item, idx) => {
            const li = document.createElement('li');
            li.textContent = item.label;
            li.addEventListener('mouseenter', () => {
                clearHighlight();
                selectedIndex = idx;
                li.style.backgroundColor = '#3498db';
            });
            li.addEventListener('mouseleave', () => {
                li.style.backgroundColor = '';
            });
            li.addEventListener('click', () => {
                window.location.href = '/artist?id=' + item.id;
            });
            box.appendChild(li);
        });

        box.style.display = 'block';
    } catch (e) {
        box.style.display = 'none';
    }
}
```
Called every time the user types in the search box:
- If the input is empty, hide the suggestions and stop
- Sends the query to `/api/search` on our Go server
- If no results come back, hide the suggestions box
- For each result, creates a list item with the label
- Highlights the item blue on hover
- Clicking an item goes directly to that artist's page
- If the request fails, hide the suggestions silently

---

## handleKey

```javascript
function handleKey(e) {
    const box = document.getElementById('suggestions');
    const items = box.querySelectorAll('li');

    if (e.key === 'ArrowDown') {
        selectedIndex = Math.min(selectedIndex + 1, items.length - 1);
        if (items[selectedIndex]) items[selectedIndex].style.backgroundColor = '#3498db';
    } else if (e.key === 'ArrowUp') {
        selectedIndex = Math.max(selectedIndex - 1, 0);
        if (items[selectedIndex]) items[selectedIndex].style.backgroundColor = '#3498db';
    } else if (e.key === 'Enter' && selectedIndex >= 0) {
        if (items[selectedIndex]) items[selectedIndex].click();
    } else if (e.key === 'Escape') {
        box.style.display = 'none';
    }
}
```
Handles keyboard navigation in the suggestions dropdown:
- `ArrowDown` — move highlight down one item
- `ArrowUp` — move highlight up one item
- `Enter` — click the currently highlighted item
- `Escape` — close the suggestions box

---

## Click Outside to Close

```javascript
document.addEventListener('click', function(e) {
    const box = document.getElementById('suggestions');
    if (!box.contains(e.target) && e.target.id !== 'q') {
        box.style.display = 'none';
    }
});
```
Listens for any click anywhere on the page. If the click was not inside the suggestions box and not on the search input, hide the suggestions.

> This is the standard pattern for closing dropdowns when the user clicks somewhere else on the page.
---
