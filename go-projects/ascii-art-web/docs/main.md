# ascii-art-web — main.go

## Package and Imports

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"ascii-art-web/art_web"
	"ascii-art-web/handlers"
)
```
This is the starting point of the program. We import:
- `fmt` — for printing messages to the terminal
- `log` — for printing fatal errors and stopping the program
- `os` — for reading the command line arguments
- `time` — for measuring how long the server takes to load
- `ascii-art-web/art_web` — our package that loads and validates banners
- `ascii-art-web/handlers` — our package that handles all HTTP routes

---

## main

```go
if len(os.Args) != 1 {
    fmt.Println("Wrong Number of Arguments")
    return
}
```
Checks that the user did not pass any extra arguments when starting the server. The only valid way to run it is:
```bash
go run .
```
If anything extra is typed after that, print an error and stop.

---

```go
start := time.Now()
```
Records the current time so we can measure how long the startup takes.

---

```go
if err := web.EnsureBanners(); err != nil {
    log.Fatalf("Banner integrity check failed: %v", err)
}
fmt.Println("✅ Banner files are present and valid.")
```
Checks that all 3 banner files exist and are valid before the server starts.
- If any banner is missing or broken — print the error and stop the program immediately
- If all banners are fine — print a success message

> We check the banners at startup so we catch problems early. If we waited until a user made a request it would be too late and they would get an error instead of art.

---

```go
elapsed := time.Since(start)
fmt.Printf("🚀 Server Finished Loading, time taken: %s\n", elapsed)
```
Calculates how long the startup took and prints it to the terminal.

> `time.Since(start)` means "how much time has passed since we recorded `start`." This is useful to know how fast the server is loading.

---

```go
handlers.Routing()
}
```
Starts the web server and keeps it running. This is the last line because `Routing()` blocks forever — it keeps listening for requests until the program is stopped.
---
