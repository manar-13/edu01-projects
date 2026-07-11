package main

import (
	"fmt"
	"log"
	"os"

	"net-cat/internal/chat"
)

const (
	defaultPort  = "8989"
	usageMessage = "[USAGE]: ./TCPChat $port\n"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Print(usageMessage)
		return
	}

	port := defaultPort
	if len(args) == 1 {
		port = args[0]
		if !chat.IsValidPort(port) {
			fmt.Print(usageMessage)
			return
		}
	}

	srv := chat.NewServer(port)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
