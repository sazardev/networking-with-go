// tcp-concurrent-server: A TCP server handling each connection concurrently with goroutines
package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintf(conn, "Echo: %s\n", text)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("Concurrent TCP server listening on :9000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}
