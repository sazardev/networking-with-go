// TCP Byte Proxy Example
// Relays raw bytes between a client and an upstream TCP server in
// both directions concurrently, with graceful shutdown on SIGINT/SIGTERM.
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os/signal"
	"syscall"
)

func proxyConn(ctx context.Context, client net.Conn, upstreamAddr string) {
	defer client.Close()

	var dialer net.Dialer
	upstream, err := dialer.DialContext(ctx, "tcp", upstreamAddr)
	if err != nil {
		fmt.Println("failed to reach upstream:", err)
		return
	}
	defer upstream.Close()

	// Copy in both directions concurrently. When either side closes,
	// both io.Copy calls return and the connections are torn down.
	done := make(chan struct{}, 2)
	go func() {
		io.Copy(upstream, client)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(client, upstream)
		done <- struct{}{}
	}()
	<-done
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()

	ln, err := net.Listen("tcp", ":9200")
	if err != nil {
		panic(err)
	}
	fmt.Println("TCP proxy listening on :9200, forwarding to :9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
		go proxyConn(ctx, conn, "localhost:9000")
	}
}
