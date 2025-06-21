package main

import (
	"fmt"
	"net"
)

// Simple TCP server that echoes messages back to the client.
func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("TCP server listening on port 9000...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	fmt.Printf("Received: %s\n", string(buf[:n]))
	conn.Write([]byte("Echo: " + string(buf[:n])))
}
