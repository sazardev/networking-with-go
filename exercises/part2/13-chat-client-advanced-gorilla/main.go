// Advanced WebSocket Chat Client (Gorilla)
// Connects to the advanced chat server, prompts for username, sends/receives JSON messages with username and timestamp.
// Usage: go run main.go
// Requires: go get github.com/gorilla/websocket

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

// Message represents the chat message format
// Matches the server's expected JSON structure
// {"user": "Ana", "msg": "¡Hola!", "time": "2025-06-22T10:01:00"}
type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
	Time string `json:"time"`
}

func main() {
	fmt.Println("=== Cliente avanzado de chat (Gorilla) ===")
	fmt.Print("Ingresa tu nombre de usuario: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = trimNewline(username)
	if username == "" {
		fmt.Println("El nombre de usuario no puede estar vacío.")
		return
	}

	// Connect to the WebSocket server
	url := "ws://localhost:8080/ws"
	fmt.Printf("Conectando a %s...\n", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error al conectar:", err)
	}
	defer conn.Close()

	// Send username as first message (handshake)
	initMsg := Message{User: username, Text: "", Time: ""}
	if err := conn.WriteJSON(initMsg); err != nil {
		log.Fatal("Error al enviar nombre de usuario:", err)
	}

	// Channel to signal shutdown
	done := make(chan struct{})

	// Goroutine: Listen for incoming messages
	go func() {
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println("\n[Desconectado del servidor]")
				close(done)
				return
			}
			// Parse and format time
			t, _ := time.Parse(time.RFC3339, msg.Time)
			fmt.Printf("[%02d:%02d] %s: %s\n", t.Hour(), t.Minute(), msg.User, msg.Text)
		}
	}()

	// Goroutine: Handle Ctrl+C (SIGINT)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("\n[Saliendo del chat]")
		conn.Close()
		os.Exit(0)
	}()

	// Main loop: Read user input and send messages
	for {
		fmt.Print("") // Prompt for input
		text, _ := reader.ReadString('\n')
		text = trimNewline(text)
		if text == "" {
			continue
		}
		msg := Message{User: username, Text: text, Time: ""}
		if err := conn.WriteJSON(msg); err != nil {
			fmt.Println("Error al enviar mensaje:", err)
			break
		}
		select {
		case <-done:
			return
		default:
		}
	}
}

// trimNewline removes trailing \r and \n from input
func trimNewline(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '\r' {
		s = s[:len(s)-1]
	}
	return s
}
