// TCP File Transfer Server Example
// Reads a size-prefixed header, then streams exactly that many
// bytes to a local file — see 14-tcp-file-transfer-client for the
// matching sender.
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	// 1. Read the 8-byte size prefix.
	var size uint64
	if err := binary.Read(conn, binary.BigEndian, &size); err != nil {
		fmt.Println("failed to read size:", err)
		return
	}

	// 2. Read the 2-byte name length, then the name itself.
	var nameLen uint16
	if err := binary.Read(conn, binary.BigEndian, &nameLen); err != nil {
		fmt.Println("failed to read name length:", err)
		return
	}
	nameBuf := make([]byte, nameLen)
	if _, err := io.ReadFull(conn, nameBuf); err != nil {
		fmt.Println("failed to read name:", err)
		return
	}
	name := string(nameBuf)

	// 3. Create the destination file and copy exactly `size` bytes into it.
	out, err := os.Create("received_" + name)
	if err != nil {
		fmt.Println("failed to create file:", err)
		return
	}
	defer out.Close()

	// io.CopyN stops after exactly `size` bytes, protecting us from a
	// misbehaving or malicious client sending more than it declared.
	written, err := io.CopyN(out, conn, int64(size))
	if err != nil && err != io.EOF {
		fmt.Println("transfer error:", err)
		return
	}
	fmt.Printf("received %q (%d bytes)\n", name, written)
}

func main() {
	ln, err := net.Listen("tcp", ":9100")
	if err != nil {
		panic(err)
	}
	fmt.Println("file server listening on :9100")
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}
