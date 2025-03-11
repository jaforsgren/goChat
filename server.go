package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients   = make(map[*Client]bool)
	clientsMu sync.Mutex
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()
	log.Println("Chat server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	client := &Client{conn: conn}
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		broadcast(msg, client)
	}

	clientsMu.Lock()
	delete(clients, client)
	clientsMu.Unlock()
	log.Printf("Client disconnected: %v", conn.RemoteAddr())
}

func broadcast(msg string, sender *Client) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		if client != sender {
			_, err := fmt.Fprintln(client.conn, msg)
			if err != nil {
				log.Printf("Failed to send message to client: %v", err)
			}
		}
	}
}
