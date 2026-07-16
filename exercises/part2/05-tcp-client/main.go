// Simple TCP Client Example
// Connects to example.com over TCP, sends a bare HTTP/1.0 request,
// and reads the full response before printing it.
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	// 1. net.Dial opens a TCP connection to example.com on port 80
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close() // Always close, even on an early return

	// 2. A deadline bounds both the write and the read below, so a
	//    stalled or silent server can't hang this goroutine forever.
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// 3. Send a simple HTTP request
	if _, err := fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n"); err != nil {
		fmt.Println("Error writing request:", err)
		os.Exit(1)
	}

	// 4. TCP is a byte stream: a single Read is not guaranteed to
	//    return the whole response, so we loop until EOF.
	reader := bufio.NewReader(conn)
	var response []byte
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		response = append(response, buf[:n]...)
		if err == io.EOF {
			break // Server closed the connection: we have it all
		}
		if err != nil {
			fmt.Println("Error reading response:", err)
			os.Exit(1)
		}
	}
	fmt.Println(string(response))
}
