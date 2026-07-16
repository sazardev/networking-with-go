// TLS Server Example
// Wraps a TCP listener with crypto/tls and greets each client over
// an encrypted channel. Generate a self-signed cert first:
//
//	openssl req -x509 -newkey rsa:2048 -nodes \
//	  -keyout server.key -out server.crt -days 365 -subj "/CN=localhost"
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
)

func main() {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		panic(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	ln, err := tls.Listen("tcp", ":9500", config)
	if err != nil {
		panic(err)
	}
	fmt.Println("TLS server listening on :9500")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		fmt.Println("shutting down, closing listener")
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return // graceful shutdown, not a real failure
			}
			fmt.Println("accept error:", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			_, err := fmt.Fprintln(c, "hello over an encrypted channel")
			if err != nil {
				fmt.Println("write to client failed:", err)
			}
		}(conn)
	}
}
