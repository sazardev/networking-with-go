// channel-broadcast: Broadcasting a message to multiple goroutines using channels
package main

import (
	"fmt"
)

func worker(id int, ch <-chan string) {
	msg := <-ch
	fmt.Printf("Worker %d received: %s\n", id, msg)
}

func main() {
	ch := make(chan string)
	for i := 1; i <= 3; i++ {
		go worker(i, ch)
	}
	for i := 1; i <= 3; i++ {
		ch <- "Broadcast message!"
	}
}
