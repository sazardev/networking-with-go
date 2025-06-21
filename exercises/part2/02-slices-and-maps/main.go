package main

import "fmt"

func main() {
	connections := []string{"client1", "client2"}
	for _, c := range connections {
		fmt.Println("Connected:", c)
	}
	ports := map[string]int{"http": 80, "https": 443}
	for name, port := range ports {
		fmt.Printf("%s => %d\n", name, port)
	}
}
