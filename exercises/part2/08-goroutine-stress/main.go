// Goroutine Stress Test: Launch 100,000 goroutines
// This demonstrates how lightweight goroutines are in Go.
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	count := 100000
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(n int) {
			defer wg.Done()
			if n == 0 {
				fmt.Println("First goroutine running!")
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("All goroutines finished!")
}
