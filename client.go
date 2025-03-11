package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to chat server")

	go readMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			break
		}
	}
}

func readMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Println(msg)
	}
	log.Println("Disconnected from server")
}
