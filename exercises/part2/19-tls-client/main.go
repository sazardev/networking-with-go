// TLS Client Example
// Connects to 19-tls-server, trusting only its self-signed
// certificate (loaded from server.crt) rather than the OS trust store.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"
)

func main() {
	caCert, err := os.ReadFile("server.crt")
	if err != nil {
		panic(err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert) {
		panic("no valid certificates found in server.crt")
	}

	config := &tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", "localhost:9500", config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		panic(fmt.Errorf("read from server failed: %w", err))
	}
	fmt.Println(string(buf[:n]))
}
