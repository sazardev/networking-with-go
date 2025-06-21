// 08-goroutines-basic/main.go
// Basic example of launching goroutines in Go.
package main

import (
	"fmt"
	"time"
)

func printMessage(msg string) {
	fmt.Println(msg)
}

func main() {
	go printMessage("Hello from goroutine 1!")
	go printMessage("Hello from goroutine 2!")
	fmt.Println("Main function running...")
	time.Sleep(500 * time.Millisecond) // Wait for goroutines to finish
	fmt.Println("Done!")
}
