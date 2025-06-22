// main.go
// Minimal WebSocket chat server using gorilla/websocket
// Broadcasts messages from any client to all connected clients.
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader upgrades HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins (for demo)
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var (
	clients   = make(map[*Client]bool) // All connected clients
	broadcast = make(chan []byte)      // Broadcast channel
	mu        sync.Mutex
)

// Handle incoming WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte)}
	mu.Lock()
	clients[client] = true
	mu.Unlock()
	defer func() {
		mu.Lock()
		delete(clients, client)
		mu.Unlock()
		conn.Close()
	}()
	// Start a goroutine to send messages to this client
	go func() {
		for msg := range client.send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}()
	// Read messages from this client and broadcast
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		broadcast <- msg
	}
}

// Broadcast messages to all clients
func handleBroadcast() {
	for msg := range broadcast {
		mu.Lock()
		for client := range clients {
			select {
			case client.send <- msg:
			default:
				// Drop client if not receiving
				close(client.send)
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleBroadcast()
	fmt.Println("Chat server running at ws://localhost:8080/ws ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
