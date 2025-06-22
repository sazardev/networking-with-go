// main.go
// WebSocket echo server using gorilla/websocket
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

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()
	for {
		// Read message from client
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", msg)
		// Echo message back to client
		if err := conn.WriteMessage(mt, msg); err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", echoHandler)
	fmt.Println("WebSocket server at ws://localhost:8080/ws ...")
	http.ListenAndServe(":8080", nil)
}
