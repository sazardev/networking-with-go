package main

import (
	"fmt"
	"time"
)

func main() {
	messages := make(chan string)
	go func() { messages <- "Ping!" }()
	fmt.Println(<-messages)

	go func() {
		fmt.Println("Handling connection in a goroutine!")
	}()
	time.Sleep(100 * time.Millisecond)
}
