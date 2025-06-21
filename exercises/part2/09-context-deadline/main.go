// 09-context-deadline/main.go
// Demonstrates context with deadline in Go.
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ch := make(chan string)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "finished work"
	}()
	select {
	case msg := <-ch:
		fmt.Println("Received:", msg)
	case <-ctx.Done():
		fmt.Println("Context deadline exceeded:", ctx.Err())
	}
}
