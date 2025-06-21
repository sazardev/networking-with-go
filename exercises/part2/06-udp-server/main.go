// Simple UDP Server Example
// Listens on UDP port 9001 and responds to each client message.
package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. net.ListenPacket creates a UDP listener
	conn, err := net.ListenPacket("udp", ":9001")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("UDP server listening on :9001")
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf) // Read message and client address
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}
		fmt.Printf("Received from %s: %s", addr, string(buf[:n]))
		// Respond to client
		conn.WriteTo([]byte("Hello from UDP server!\n"), addr)
	}
}
