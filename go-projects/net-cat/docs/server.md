# net-cat — internal/chat/server.go

## Constants

```go
const (
	maxConnections = 10
)
```
The maximum number of clients allowed in the chat at the same time.

---

## Server Struct

```go
type Server struct {
	port string

	clientsMu sync.Mutex
	clients   map[net.Conn]*Client

	historyMu sync.Mutex
	history   []string

	joinCh      chan *Client
	leaveCh     chan *Client
	broadcastCh chan string

	namesMu sync.Mutex
	names   map[string]bool
}
```
Holds everything the server needs to run:
- `port` — the port the server listens on
- `clientsMu` + `clients` — a mutex and a map of all currently connected clients keyed by their connection
- `historyMu` + `history` — a mutex and a list of all messages sent so far — given to new clients when they join
- `joinCh` — a channel that receives a client when they finish setting up and are ready to join
- `leaveCh` — a channel that receives a client when they disconnect
- `broadcastCh` — a channel that receives messages to be sent to all clients
- `namesMu` + `names` — a mutex and a map of all taken names to prevent duplicates

> Every shared piece of data has its own mutex. This is important because many goroutines run at the same time — one per client. Without mutexes they would corrupt each other's data.

---

## NewServer

```go
func NewServer(port string) *Server {
	return &Server{
		port:        port,
		clients:     make(map[net.Conn]*Client),
		history:     make([]string, 0, 128),
		joinCh:      make(chan *Client),
		leaveCh:     make(chan *Client),
		broadcastCh: make(chan string, 256),
		names:       make(map[string]bool),
	}
}
```
Creates and returns a new server with all fields initialized.
- History starts with capacity 128 — pre-allocates space for 128 messages before needing to grow
- `broadcastCh` has a buffer of 256 — can hold 256 messages before blocking
- `joinCh` and `leaveCh` are unbuffered — the sender waits until the loop receives them

---

## Run

```go
func (s *Server) Run() error {
	ln, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Printf("Listening on the port :%s\n", s.port)

	go s.loop()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		s.clientsMu.Lock()
		if len(s.clients) >= maxConnections {
			s.clientsMu.Unlock()
			_, _ = conn.Write([]byte("Chat is full (max 10). Try later.\n"))
			_ = conn.Close()
			continue
		}
		s.clientsMu.Unlock()

		go s.handleConn(conn)
	}
}
```
Starts the TCP server and accepts incoming connections:
- Creates a TCP listener on the given port
- `defer ln.Close()` — closes the listener when the function ends
- Prints the listening message
- Starts the event loop in a separate goroutine
- Loops forever accepting new connections
- If accepting fails, log the error and continue — don't crash
- If the chat is full, tell the client and close their connection
- Otherwise start a new goroutine to handle that connection

> `go s.handleConn(conn)` means each client gets their own goroutine. The server can handle many clients at the same time because they all run independently.

---

## loop

```go
func (s *Server) loop() {
	for {
		select {
		case c := <-s.joinCh:
```
The central event loop — runs in its own goroutine and handles all join, leave, and broadcast events one at a time.

> Using one loop for all events means we never have race conditions between joining, leaving, and broadcasting. Everything happens in order.

---

```go
			s.clientsMu.Lock()
			s.clients[c.conn] = c
			s.clientsMu.Unlock()

			s.namesMu.Lock()
			s.names[c.name] = true
			s.namesMu.Unlock()

			for _, h := range s.copyHistory() {
				select {
				case c.out <- h:
				default:
				}
			}

			s.broadcast(formatSystem(c.name + " has joined our chat..."))
```
When a client joins:
- Add them to the clients map
- Mark their name as taken
- Send them all previous messages so they can catch up
- Broadcast a join notification to all other clients

> `select { case c.out <- h: default: }` tries to send the history message but skips it if the client's channel is full. This prevents the server from blocking if a slow client can't keep up.

---

```go
		case c := <-s.leaveCh:
			s.clientsMu.Lock()
			if client, ok := s.clients[c.conn]; ok {
				delete(s.clients, c.conn)
				s.namesMu.Lock()
				delete(s.names, client.name)
				s.namesMu.Unlock()
				close(c.out)
			}
			s.clientsMu.Unlock()
			s.broadcast(formatSystem(c.name + " has left our chat..."))
```
When a client leaves:
- Remove them from the clients map
- Free their name so someone else can use it
- Close their `out` channel which stops their writer goroutine
- Broadcast a leave notification to all remaining clients

---

```go
		case msg := <-s.broadcastCh:
			s.appendHistory(msg)
			s.clientsMu.Lock()
			for _, cl := range s.clients {
				select {
				case cl.out <- msg:
				default:
				}
			}
			s.clientsMu.Unlock()
		}
	}
}
```
When a message arrives to broadcast:
- Save it to history so new clients can see it
- Loop through every connected client and send the message to their `out` channel
- Skip any client whose channel is full

---

## handleConn

```go
func (s *Server) handleConn(conn net.Conn) {
	if _, err := conn.Write([]byte(linuxLogo)); err != nil {
		_ = conn.Close()
		return
	}
```
Handles one client from connection to disconnection:
- Sends the welcome logo immediately when connected
- If sending fails the client already disconnected — close and return

---

```go
	for {
		name, err = getClientName(conn)
		if err != nil {
			_ = conn.Close()
			return
		}

		s.namesMu.Lock()
		_, exists := s.names[name]
		s.namesMu.Unlock()

		if !exists {
			break
		}

		if _, err := conn.Write([]byte("Name already taken. Please choose another: ")); err != nil {
			_ = conn.Close()
			return
		}
	}
```
Keeps asking for a name until the client gives one that is not already taken:
- Calls `getClientName` to read the name from the client
- Checks if the name is already in the names map
- If it is taken, tell the client and ask again
- If it is free, break out of the loop and continue

---

```go
	client := &Client{
		conn: conn,
		name: name,
		out:  make(chan string, 64),
	}

	go client.writer()

	s.joinCh <- client

	client.reader(func(text string, from *Client) {
		s.broadcast(formatChat(from.name, text))
	})

	s.leaveCh <- client
	_ = client.conn.Close()
}
```
Once the client has a valid name:
- Creates the `Client` struct with a buffered output channel of 64 messages
- Starts the writer goroutine to send messages to this client
- Sends the client to the join channel — the loop adds them to the chat
- Starts reading messages from the client — this blocks until the client disconnects
- When `reader` returns (client disconnected), sends to the leave channel and closes the connection

---

## broadcast, appendHistory, copyHistory

```go
func (s *Server) broadcast(msg string) {
	s.broadcastCh <- msg
}
```
A simple helper that sends a message to the broadcast channel.

---

```go
func (s *Server) appendHistory(line string) {
	s.historyMu.Lock()
	defer s.historyMu.Unlock()
	s.history = append(s.history, line)
	if len(s.history) > 1000 {
		s.history = s.history[len(s.history)-1000:]
	}
}
```
Adds a message to the history list safely using a mutex. If history grows beyond 1000 messages, keeps only the most recent 1000.

> Limiting history to 1000 messages prevents the server from using too much memory over time.

---

```go
func (s *Server) copyHistory() []string {
	s.historyMu.Lock()
	defer s.historyMu.Unlock()
	cp := make([]string, len(s.history))
	copy(cp, s.history)
	return cp
}
```
Returns a safe copy of the history list. Uses a mutex and copies the slice so the caller can read it without holding the lock.

> Returning a copy instead of the original slice is important. If we returned the original, another goroutine could modify it while the new client is reading it — causing a data race.
---
