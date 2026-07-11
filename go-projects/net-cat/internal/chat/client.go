package chat

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"strings"
)

type Client struct {
	conn net.Conn
	name string
	out  chan string
}

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

func (c *Client) writer() {
	for msg := range c.out {
		if _, err := c.conn.Write([]byte(msg + "\n")); err != nil {
			_ = c.conn.Close()
			return
		}
	}
}

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
