// TCP File Transfer Client Example
// Sends a local file to 14-tcp-file-transfer-server using a small
// size-prefixed header protocol.
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

// dialWithRetry attempts to connect up to attempts times, backing off
// between failures. A bare net.Dial in production code is often too
// fragile: a server still starting up, a brief network blip, or a load
// balancer mid-failover can all cause a single dial to fail even though
// the destination is reachable moments later.
func dialWithRetry(
	addr string, attempts int, backoff time.Duration,
) (net.Conn, error) {
	dialer := net.Dialer{Timeout: 5 * time.Second}

	var lastErr error
	for i := 0; i < attempts; i++ {
		conn, err := dialer.Dial("tcp", addr)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		time.Sleep(backoff * time.Duration(i+1))
	}
	return nil, fmt.Errorf("after %d attempts: %w", attempts, lastErr)
}

func sendFile(path, addr string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	conn, err := dialWithRetry(addr, 3, 500*time.Millisecond)
	if err != nil {
		return err
	}
	defer conn.Close()

	name := filepath.Base(path)

	// Write the header: size, then name length, then name.
	size := uint64(info.Size())
	if err := binary.Write(conn, binary.BigEndian, size); err != nil {
		return err
	}
	nameLen := uint16(len(name))
	if err := binary.Write(conn, binary.BigEndian, nameLen); err != nil {
		return err
	}
	if _, err := conn.Write([]byte(name)); err != nil {
		return err
	}

	// io.Copy streams the file in chunks — it never loads the whole
	// file into memory, so this works the same for a 1 KB file or a 10 GB one.
	sent, err := io.Copy(conn, f)
	if err != nil {
		return err
	}
	fmt.Printf("sent %d bytes\n", sent)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: 14-tcp-file-transfer-client <path>")
		os.Exit(1)
	}
	if err := sendFile(os.Args[1], "localhost:9100"); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
