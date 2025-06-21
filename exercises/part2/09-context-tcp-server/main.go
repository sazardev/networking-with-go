// 09-context-tcp-server/main.go
// TCP server that supports cancellation via context.
package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleConn(ctx context.Context, c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)
	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(c, "Server shutting down.")
			return
		default:
			c.SetReadDeadline(time.Now().Add(1 * time.Second))
			msg, err := reader.ReadString('\n')
			if err != nil {
				continue
			}
			fmt.Fprintf(c, "Echo: %s", msg)
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on :9090 (Ctrl+C to stop)...")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	for {
		conn, err := ln.Accept()
		if err != nil {
			break
		}
		go handleConn(ctx, conn)
		select {
		case <-ctx.Done():
			fmt.Println("Server context cancelled, shutting down...")
			ln.Close()
			return
		default:
		}
	}
}
