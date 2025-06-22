package main

import (
	"crypto/tls"
	"fmt"
	"log"
)

// Example TLS server (requires cert.pem and key.pem in the same directory)
func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", ":8443", config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TLS server listening on port 8443...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go func(c *tls.Conn) {
			defer c.Close()
			c.Write([]byte("Hello, secure world!\n"))
		}(conn.(*tls.Conn))
	}
}
