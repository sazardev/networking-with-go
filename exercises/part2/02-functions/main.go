package main

import "fmt"

func greet(name string) string {
	return "Welcome, " + name + "!"
}

func main() {
	fmt.Println(greet("Gopher"))
}
