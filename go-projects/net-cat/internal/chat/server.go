package chat

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	maxConnections = 10
)

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

func (s *Server) loop() {
	for {
		select {
		case c := <-s.joinCh:
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

func (s *Server) handleConn(conn net.Conn) {

	if _, err := conn.Write([]byte(linuxLogo)); err != nil {
		_ = conn.Close()
		return
	}

	var name string
	var err error

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

func (s *Server) broadcast(msg string) {
	s.broadcastCh <- msg
}

func (s *Server) appendHistory(line string) {
	s.historyMu.Lock()
	defer s.historyMu.Unlock()
	s.history = append(s.history, line)
	if len(s.history) > 1000 {
		s.history = s.history[len(s.history)-1000:]
	}
}

func (s *Server) copyHistory() []string {
	s.historyMu.Lock()
	defer s.historyMu.Unlock()
	cp := make([]string, len(s.history))
	copy(cp, s.history)
	return cp
}
