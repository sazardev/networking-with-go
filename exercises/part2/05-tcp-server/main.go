// Simple TCP Server Example
// Listens on port 8080 and greets each client, shutting down
// cleanly on Ctrl+C.
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
)

func main() {
	// 1. net.Listen creates a TCP listener on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Server listening on :8080")

	// 2. Cancel this context on Ctrl+C to shut down gracefully
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down, closing listener...")
		ln.Close() // Unblocks Accept below with an error
	}()

	for {
		// 3. Accept blocks until a new connection arrives
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return // We closed the listener on purpose
			default:
				fmt.Println("Error accepting connection:", err)
				continue
			}
		}
		// 4. Handle each connection in its own goroutine
		go func(c net.Conn) {
			defer c.Close() // Runs even if this goroutine panics
			fmt.Fprintln(c, "Hello from the Go TCP server!")
		}(conn)
	}
}
