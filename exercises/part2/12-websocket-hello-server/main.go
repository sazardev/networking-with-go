// main.go
// Native WebSocket server that greets each client with their IP address.
// No third-party packages. For learning purposes only.
package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
)

const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// handleWS handles the WebSocket handshake and sends a greeting with the client's IP.
func handleWS(w http.ResponseWriter, r *http.Request) {
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
	// Write handshake response
	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", accept)
	w.WriteHeader(http.StatusSwitchingProtocols)
	// Hijack the connection
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	conn, buf, err := hj.Hijack()
	if err != nil {
		return
	}
	defer conn.Close()
	// Get client IP
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	greeting := fmt.Sprintf("Hello! Your IP is %s", ip)
	// Send greeting as a WebSocket text frame
	sendTextFrame(buf, greeting)
	buf.Flush()
	// Optionally, you could now enter a loop to echo or handle more frames
}

func computeAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key + wsGUID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// sendTextFrame writes a single WebSocket text frame with the given message.
func sendTextFrame(w io.Writer, msg string) {
	payload := []byte(msg)
	frame := []byte{0x81} // FIN=1, opcode=1 (text)
	if len(payload) < 126 {
		frame = append(frame, byte(len(payload)))
	} else if len(payload) < 65536 {
		frame = append(frame, 126, byte(len(payload)>>8), byte(len(payload)))
	} else {
		// Not handling >64K payloads for simplicity
		return
	}
	frame = append(frame, payload...)
	w.Write(frame)
}

func main() {
	http.HandleFunc("/ws", handleWS)
	fmt.Println("WebSocket Hello server running at ws://localhost:8080/ws")
	http.ListenAndServe(":8080", nil)
}
