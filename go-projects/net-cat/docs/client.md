# net-cat — internal/chat/client.go

## Client Struct

```go
type Client struct {
	conn net.Conn
	name string
	out  chan string
}
```
Represents one connected client in the chat.
- `conn` — the TCP connection to this client
- `name` — the name the client chose when they joined
- `out` — a channel that holds messages waiting to be sent to this client

> The `out` channel is the key design here. Instead of writing directly to the connection from multiple places, everything goes through this one channel. This prevents two goroutines from writing to the same connection at the same time.

---

## getClientName

```go
func getClientName(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	for {
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		name = strings.TrimSpace(name)
		if name != "" {
			return name, nil
		}
		if _, err := conn.Write([]byte("[ENTER YOUR NAME]: ")); err != nil {
			return "", err
		}
	}
}
```
Reads the name the client types and keeps asking until a non-empty name is given.
- Creates a buffered reader on the connection
- Waits for the client to type something and press Enter
- If the connection drops, return an error
- Removes any extra spaces or newlines from the name
- If the name is not empty, return it
- If the name is empty, ask again by sending `[ENTER YOUR NAME]:` back to the client
- Loops until a valid name is received

> `bufio.NewReader` is used because TCP data arrives in a stream. The buffered reader waits until it sees a newline `\n` before returning the full line — this is how we know the user finished typing.

---

## writer

```go
func (c *Client) writer() {
	for msg := range c.out {
		if _, err := c.conn.Write([]byte(msg + "\n")); err != nil {
			_ = c.conn.Close()
			return
		}
	}
}
```
Runs in its own goroutine and sends messages to this client one by one.
- Loops through every message that arrives in the `out` channel
- Writes each message to the client's TCP connection followed by a newline
- If writing fails (client disconnected), close the connection and stop
- When the `out` channel is closed, the loop ends automatically

> `for msg := range c.out` is a Go pattern for reading from a channel until it is closed. This goroutine stays alive the whole time the client is connected, always ready to send the next message.

---

## reader

```go
func (c *Client) reader(onMessage func(text string, from *Client)) {
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println("read error:", err)
			}
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		onMessage(line, c)
	}
}
```
Runs and waits for the client to type messages, then passes each one to a callback function.
- Creates a buffered reader on the connection
- Waits for a full line from the client
- If the client disconnects (EOF), stop quietly without logging
- If any other error happens, log it and stop
- Removes extra spaces and newlines from the message
- If the message is empty after trimming, skip it — we do not broadcast empty messages
- Calls `onMessage` with the text and the client who sent it

> `onMessage` is a function passed in from the server. When the reader calls it, the server broadcasts the message to all other clients. This keeps the reading logic separate from the broadcasting logic.
---