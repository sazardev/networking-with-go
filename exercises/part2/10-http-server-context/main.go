// 10-http-server-context/main.go
// Servidor HTTP usando context para cancelar peticiones largas
package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("Petición recibida, procesando...")

	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintln(w, "¡Procesamiento completado!")
	case <-ctx.Done():
		fmt.Fprintln(w, "Petición cancelada por el cliente.")
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Servidor con context en http://localhost:8083 ...")
	http.ListenAndServe(":8083", nil)
}
