// Logging TCP Server Example
// Listens on port 8081 and logs every accepted connection with
// the standard log package instead of fmt.Println.
package main

import (
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
	log.Println("Server listening on :8081")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("Accepted connection from %v", conn.RemoteAddr())
		conn.Close()
	}
}
