// 10-http-server-basic/main.go
// Servidor HTTP básico en Go
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "¡Hola, mundo! Servidor HTTP básico en Go.")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Servidor escuchando en http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}
