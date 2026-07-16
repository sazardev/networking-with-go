// TCP Echo Server Example
// Returns exactly what each client sends, disconnecting idle
// clients after 30 seconds of silence.
package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	ln, _ := net.Listen("tcp", ":9000")
	fmt.Println("Echo server on :9000")
	for {
		conn, _ := ln.Accept()
		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 1024)
			for {
				// Reset the deadline on every loop: an idle client
				// gets disconnected after 30s of silence instead of
				// holding the goroutine open forever.
				c.SetReadDeadline(time.Now().Add(30 * time.Second))
				n, err := c.Read(buf) // Read data from the client
				if n > 0 {
					// Write can be short: check its error, since
					// net.Conn.Write only returns n < len(p) together
					// with a non-nil error.
					if _, werr := c.Write(buf[:n]); werr != nil {
						fmt.Println("Write error:", werr)
						return
					}
				}
				if err != nil {
					if !errors.Is(err, io.EOF) {
						fmt.Println("Read error:", err)
					}
					return // Client closed, or the deadline fired
				}
			}
		}(conn)
	}
}
