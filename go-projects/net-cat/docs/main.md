# net-cat — main.go

## Constants

```go
const (
	defaultPort  = "8989"
	usageMessage = "[USAGE]: ./TCPChat $port\n"
)
```
Two constants used throughout the main function:
- `defaultPort` — if no port is given, the server starts on port 8989
- `usageMessage` — the error message shown when the wrong arguments are given

---

## main

```go
args := os.Args[1:]

if len(args) > 1 {
	fmt.Print(usageMessage)
	return
}
```
Reads all command line arguments after the program name.
- If more than one argument is given, print the usage message and stop
- Only zero or one argument is valid

---

```go
port := defaultPort
if len(args) == 1 {
	port = args[0]
	if !chat.IsValidPort(port) {
		fmt.Print(usageMessage)
		return
	}
}
```
Decides which port to use:
- Starts with the default port `8989`
- If one argument was given, use it as the port
- Validates the port using `chat.IsValidPort` — if it is not a valid port number, print the usage message and stop

---

```go
srv := chat.NewServer(port)
if err := srv.Run(); err != nil {
	log.Fatal(err)
}
```
Creates the server and starts it:
- `chat.NewServer(port)` — creates a new server with the chosen port
- `srv.Run()` — starts listening for connections and blocks forever
- If the server fails to start (for example the port is already in use), print the error and stop

> `main.go` is intentionally simple. It only handles arguments and starts the server. All the real work happens in the `chat` package.
---
