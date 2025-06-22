// json_client.go
// HTTP client that fetches users and posts a new user to the JSON server.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// User represents a user with ID, Name, and Email fields.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Fetch users (GET)
	resp, err := http.Get("http://localhost:8080/users") // Send GET request
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil { // Decode JSON array
		panic(err)
	}
	fmt.Println("Users from server:")
	for _, u := range users {
		fmt.Printf("- ID: %d, Name: %s, Email: %s\n", u.ID, u.Name, u.Email)
	}

	// Add a new user (POST)
	newUser := User{Name: "Charlie", Email: "charlie@example.com"}
	data, _ := json.Marshal(newUser) // Encode new user as JSON
	resp2, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewBuffer(data)) // Send POST
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	var created User
	body, _ := io.ReadAll(resp2.Body)
	json.Unmarshal(body, &created) // Decode response
	fmt.Println("Created user:", created)
}
