# groupie-tracker — fetch.go

## fetchJSON

```go
func fetchJSON[T any](url string, target *T) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to GET from %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v", err)
	}

	return json.Unmarshal(body, target)
}
```
A generic function that fetches JSON data from a URL and stores it in any type of variable.
- `[T any]` means this function works with any data type — Artist, RelationData, LocationData, etc.
- Sends a GET request to the URL
- If the request fails, return an error with the URL in the message so we know which one failed
- `defer resp.Body.Close()` — closes the response body automatically when the function finishes
- Reads the full response body into `body`
- Converts the JSON bytes into the Go data type and stores it in `target`

> `json.Unmarshal` is what converts raw JSON text into Go structs. For example it turns `{"id":1,"name":"Queen"}` into an `Artist` struct with `ID=1` and `Name="Queen"`.

---

## fetchJSONCtx

```go
func fetchJSONCtx[T any](ctx context.Context, url string, target *T) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to GET from %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v", err)
	}

	return json.Unmarshal(body, target)
}
```
The same as `fetchJSON` but with a context — used for the startup data fetching with a timeout.
- Takes a `ctx` (context) as an extra parameter
- `http.NewRequestWithContext` creates the request with the context attached — if the context times out or is cancelled, the request stops automatically
- Creates a new HTTP client and sends the request
- Everything else is the same as `fetchJSON`

> The difference between the two functions is the context. `fetchJSONCtx` is used at startup where we set a 12-second timeout — if the API does not respond in time the server stops instead of hanging forever. `fetchJSON` is simpler and used where no timeout is needed.
---
