// json_server.go
// HTTP server that serves a list of users as JSON and receives new users via POST.
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users = []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}
	mu sync.Mutex
)

// getUsers handles GET /users and returns the list as JSON.
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(users)
}

// addUser handles POST /users and adds a new user from JSON body.
func addUser(w http.ResponseWriter, r *http.Request) {
	var u User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	u.ID = len(users) + 1
	users = append(users, u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsers(w, r)
		} else if r.Method == http.MethodPost {
			addUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8080", nil)
}
