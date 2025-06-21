// 09-context-cancel-basic/main.go
// Demonstrates basic context cancellation in Go.
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Goroutine: context cancelled!")
		}
	}()
	fmt.Println("Main: sleeping 1s...")
	time.Sleep(1 * time.Second)
	fmt.Println("Main: calling cancel()")
	cancel()
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Main: done.")
}
