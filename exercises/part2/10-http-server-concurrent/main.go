// 10-http-server-concurrent/main.go
// Servidor HTTP concurrente en Go
package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var counter int64

func handler(w http.ResponseWriter, r *http.Request) {
	n := atomic.AddInt64(&counter, 1)
	fmt.Fprintf(w, "Petición #%d atendida concurrentemente\n", n)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Servidor concurrente en http://localhost:8082 ...")
	http.ListenAndServe(":8082", nil)
}
