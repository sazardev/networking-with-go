// Advanced WebSocket Chat Server with Username and Timestamp (Gorilla)
// Run: go run exercises/part2/13-chat-server-advanced-gorilla/main.go
// Requires: go get github.com/gorilla/websocket

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
	Time string `json:"time"`
}

type Client struct {
	conn *websocket.Conn
	user string
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan Message)
	mutex     sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// First message must be the username
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Username read error:", err)
		return
	}
	var join struct{ User string }
	if err := json.Unmarshal(msg, &join); err != nil || join.User == "" {
		log.Println("Invalid join message")
		return
	}
	client := &Client{conn: conn, user: join.User}

	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	log.Printf("%s joined the chat", client.user)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}
		m.User = client.user
		m.Time = time.Now().Format(time.RFC3339)
		broadcast <- m
	}

	mutex.Lock()
	delete(clients, client)
	mutex.Unlock()
	log.Printf("%s left the chat", client.user)
}

func handleMessages() {
	for {
		msg := <-broadcast
		msgBytes, _ := json.Marshal(msg)
		mutex.Lock()
		for client := range clients {
			client.conn.WriteMessage(websocket.TextMessage, msgBytes)
		}
		mutex.Unlock()
	}
}
