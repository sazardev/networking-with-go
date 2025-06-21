package main

import "fmt"

func main() {
	for i := 1; i <= 3; i++ {
		fmt.Println("Packet", i)
	}
	port := 8080
	if port == 8080 {
		fmt.Println("Standard HTTP port!")
	}
}
