// main.go
// HTTP server that serves a list of users as JSON and receives new users via POST.
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

// User represents a user with ID, Name, and Email fields.
type User struct {
	ID    int    `json:"id"`    // Unique identifier for the user
	Name  string `json:"name"`  // User's name
	Email string `json:"email"` // User's email address
}

var (
	// users holds the in-memory list of users.
	users = []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}
	// mu is a mutex to protect concurrent access to users slice.
	mu sync.Mutex
)

// getUsers handles GET /users and returns the list as JSON.
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set response type
	mu.Lock()                                          // Lock to safely read users
	defer mu.Unlock()
	json.NewEncoder(w).Encode(users) // Encode users slice as JSON
}

// addUser handles POST /users and adds a new user from JSON body.
func addUser(w http.ResponseWriter, r *http.Request) {
	var u User
	body, err := io.ReadAll(r.Body) // Read request body
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &u); err != nil { // Parse JSON into User
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	mu.Lock() // Lock to safely modify users
	defer mu.Unlock()
	u.ID = len(users) + 1    // Assign a new ID
	users = append(users, u) // Add to list
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u) // Respond with the created user
}

func main() {
	// Route /users to GET and POST handlers
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsers(w, r)
		} else if r.Method == http.MethodPost {
			addUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", nil)
}
