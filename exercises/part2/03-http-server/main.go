package main

import (
	"fmt"
	"net/http"
)

// Simple HTTP server that responds with a greeting.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go HTTP server!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
