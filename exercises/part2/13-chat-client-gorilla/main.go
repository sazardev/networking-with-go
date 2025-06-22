// main.go
// Minimal Go WebSocket chat client using gorilla/websocket
// Connects to a chat server, sends user input, and prints all received messages.
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

func main() {
	// Connect to the WebSocket server
	url := "ws://localhost:8080/ws"
	fmt.Println("Connecting to", url)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer c.Close()

	// Channel to handle OS interrupts (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Goroutine to read messages from the server
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			fmt.Println("[Server]:", string(msg))
		}
	}()

	// Read user input and send to server
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type messages and press Enter to send. Ctrl+C to exit.")
	for {
		select {
		case <-interrupt:
			fmt.Println("Interrupted, closing connection...")
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
			time.Sleep(time.Second)
			return
		default:
			if scanner.Scan() {
				msg := scanner.Text()
				if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					fmt.Println("Write error:", err)
					return
				}
			}
		}
	}
}
