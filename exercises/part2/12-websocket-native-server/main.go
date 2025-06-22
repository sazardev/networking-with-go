// main.go
// Minimal WebSocket echo server using only Go's standard library (no gorilla/websocket)
// Note: The Go standard library does not provide a high-level WebSocket API, so we must handle the handshake and frames manually.
// This is for educational purposes and not recommended for production.
package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

// GUID as per RFC 6455
const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func computeAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key + wsGUID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Connection") != "Upgrade" || r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		http.Error(w, "Missing Sec-WebSocket-Key", http.StatusBadRequest)
		return
	}
	accept := computeAcceptKey(key)
	// Hijack the connection
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	conn, buf, err := hj.Hijack()
	if err != nil {
		http.Error(w, "Hijack failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	// Write the WebSocket handshake response
	fmt.Fprintf(conn, "HTTP/1.1 101 Switching Protocols\r\n")
	fmt.Fprintf(conn, "Upgrade: websocket\r\n")
	fmt.Fprintf(conn, "Connection: Upgrade\r\n")
	fmt.Fprintf(conn, "Sec-WebSocket-Accept: %s\r\n\r\n", accept)
	// Echo loop: read a frame, send it back
	for {
		// Read the first two bytes (FIN/Opcode, Mask/PayloadLen)
		head := make([]byte, 2)
		if _, err := io.ReadFull(buf, head); err != nil {
			break
		}
		fin := head[0]&0x80 != 0
		opcode := head[0] & 0x0F
		masked := head[1]&0x80 != 0
		payloadLen := int(head[1] & 0x7F)
		if payloadLen == 126 {
			// 2 bytes extended
			ext := make([]byte, 2)
			io.ReadFull(buf, ext)
			payloadLen = int(ext[0])<<8 | int(ext[1])
		} else if payloadLen == 127 {
			// 8 bytes extended (not handled for simplicity)
			break
		}
		var maskKey [4]byte
		if masked {
			io.ReadFull(buf, maskKey[:])
		}
		payload := make([]byte, payloadLen)
		io.ReadFull(buf, payload)
		if masked {
			for i := 0; i < payloadLen; i++ {
				payload[i] ^= maskKey[i%4]
			}
		}
		if opcode == 0x8 { // Close frame
			break
		}
		if opcode == 0x1 && fin { // Text frame
			// Echo back
			resp := []byte{0x81, byte(len(payload))}
			conn.Write(resp)
			conn.Write(payload)
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Native WebSocket echo server at ws://localhost:8082/ws")
	http.ListenAndServe(":8082", nil)
}
