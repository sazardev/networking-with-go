// main.go
// Minimal WebSocket client using only Go's standard library (no gorilla/websocket)
// Note: This is for educational purposes and only supports simple text frames.
package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// Generate a random key for the handshake
	keyBytes := make([]byte, 16)
	rand.Read(keyBytes)
	key := base64.StdEncoding.EncodeToString(keyBytes)
	// Send the handshake request
	req := "GET /ws HTTP/1.1\r\n" +
		"Host: localhost:8082\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: " + key + "\r\n" +
		"Sec-WebSocket-Version: 13\r\n\r\n"
	conn.Write([]byte(req))
	// Read handshake response
	resp, _ := bufio.NewReader(conn).ReadString('\n')
	if !strings.Contains(resp, "101") {
		fmt.Println("Handshake failed:", resp)
		return
	}
	// Skip headers
	for {
		line, _ := bufio.NewReader(conn).ReadString('\n')
		if line == "\r\n" {
			break
		}
	}
	fmt.Println("Connected! Type messages to send, Ctrl+C to quit.")
	// Send a text frame and read echo
	for {
		fmt.Print("> ")
		var msg string
		fmt.Scanln(&msg)
		if msg == "" {
			continue
		}
		// Send text frame (FIN=1, opcode=1)
		frame := []byte{0x81, byte(len(msg))}
		frame = append(frame, []byte(msg)...)
		conn.Write(frame)
		// Read echo
		head := make([]byte, 2)
		io.ReadFull(conn, head)
		payloadLen := int(head[1] & 0x7F)
		payload := make([]byte, payloadLen)
		io.ReadFull(conn, payload)
		fmt.Println("Echo:", string(payload))
	}
}
