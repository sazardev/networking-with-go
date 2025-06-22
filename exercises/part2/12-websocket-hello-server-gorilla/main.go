// main.go
// WebSocket hello server using gorilla/websocket
// Greets each client with their IP address upon connection.
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader upgrades HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins (for demo)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Get the client's IP address
	clientIP := r.RemoteAddr
	// Send a greeting message with the client's IP
	greeting := fmt.Sprintf("Hello! Your IP is %s", clientIP)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(greeting)); err != nil {
		fmt.Println("Write error:", err)
		return
	}

	// Optionally, keep the connection open for further messages (echo loop)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Printf("Received from %s: %s\n", clientIP, msg)
		// Echo the message back
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", helloHandler)
	fmt.Println("WebSocket hello server at ws://localhost:8080/ws ...")
	http.ListenAndServe(":8080", nil)
}
