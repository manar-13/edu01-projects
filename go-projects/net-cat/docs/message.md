# net-cat — internal/chat/message.go

## linuxLogo

```go
const linuxLogo = "Welcome to TCP-Chat!\n" +
	"         _nnnn_\n" +
	"        dGGGGMMb\n" +
	...
	"[ENTER YOUR NAME]: "
```
A constant string that holds the welcome message and ASCII art logo sent to every client when they first connect.
- Built by joining multiple strings together with `+`
- Each line ends with `\n` to create a new line in the terminal
- The last line is `[ENTER YOUR NAME]: ` which prompts the client to type their name immediately after the logo

> This is stored as a constant because it never changes. Every client that connects sees exactly the same welcome message.

---

## formatChat

```go
func formatChat(name, text string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:%s", t, name, text)
}
```
Formats a chat message with a timestamp and the sender's name.
- Gets the current time and formats it as `"2006-01-02 15:04:05"`
- Returns a string in the format `[2020-01-20 15:48:41][Manar]:hello`

> Go uses a specific reference time `2006-01-02 15:04:05` for formatting dates. This exact time is Go's birthday — January 2, 2006 at 15:04:05. You use it as a template to define how you want the date to look.

---

## formatSystem

```go
func formatSystem(text string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][Server]:%s", t, text)
}
```
Formats a system message — used for join and leave notifications.
- Same timestamp format as `formatChat`
- Uses `[Server]` instead of a client name to show this message is from the server
- Returns a string like `[2020-01-20 16:04:10][Server]:Lee has joined our chat...`

> The difference between `formatChat` and `formatSystem` is the sender label. Chat messages show the client name, system messages show `Server`. This makes it clear to everyone in the chat whether a message was sent by a person or by the server itself.
---