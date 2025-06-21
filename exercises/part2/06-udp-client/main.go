// Simple UDP Client Example
// Sends a message to a UDP server and prints the response.
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 1. net.Dial creates a UDP connection (no handshake)
	conn, err := net.Dial("udp", "localhost:9001")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	// 2. Send a message
	fmt.Fprintf(conn, "Hello UDP server!\n")
	// 3. Read the response
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("Response:", string(buf[:n]))
	conn.Close()
}
