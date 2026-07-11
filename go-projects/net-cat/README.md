# TCP Chat (Net-Cat Clone)

A simple TCP chat application written in Go, inspired by the `netcat` (`nc`) utility.  
Supports multiple clients in a group chat with join/leave notifications, message history, and timestamps.

---

## Features

- ✅ Server–Client architecture over **TCP**
- ✅ Up to **10 concurrent clients**
- ✅ Clients must provide a **non-empty name**
- ✅ **Message format**: `[YYYY-MM-DD HH:MM:SS][name]:message`
- ✅ Broadcasts **join** and **leave** events
- ✅ New clients receive the **full chat history**
- ✅ Ignores and does not broadcast empty messages
- ✅ Default port **8989**, or custom port via CLI
- ✅ Graceful handling of client disconnects

---

## Usage

### Start the server

```bash
# Run on default port 8989
go run .

# Run on custom port
go run . 2525

```
## Connect as a client
- Use nc (netcat) or telnet from another terminal:

```bash
nc localhost 8989
```
- On connect you'll see:
```bash
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
```

- Example Chat:
```bash
[2025-09-30 14:15:01][Alice]:Hello everyone!
[2025-09-30 14:15:10][Bob]:Hi Alice!
Alice has joined our chat...
Bob has joined our chat...
[2025-09-30 14:15:45][Alice]:How are you?
[2025-09-30 14:15:53][Bob]:Great, thanks!
Bob has left our chat...
```
## Project Structure 
```bash
.
├── go.mod
├── main.go              # Entry point
└── internal/chat
    ├── client.go        # Client connection handling
    ├── server.go        # Server hub & broadcast
    ├── message.go       # Formatting & ASCII banner
    └── validate.go      # Port validation
```
## Authors:
* Manar Mohamed (manmohamed)



