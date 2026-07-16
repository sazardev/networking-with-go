// Testing HTTP Handlers Example
// A couple of minimal handlers exercised by main_test.go using
// net/http/httptest, without ever binding a real port for the
// unit test and a random free port for the end-to-end one.
package main

import (
	"io"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if _, err := io.Copy(io.Discard, r.Body); err != nil {
		http.Error(w, "read failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8084", nil)
}
