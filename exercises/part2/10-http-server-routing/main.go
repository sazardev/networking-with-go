// 10-http-server-routing/main.go
// Servidor HTTP con rutas personalizadas en Go
package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "¡Hola desde /hello!")
}

func byeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "¡Adiós desde /bye!")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/bye", byeHandler)
	fmt.Println("Servidor escuchando en http://localhost:8081 ...")
	http.ListenAndServe(":8081", nil)
}
