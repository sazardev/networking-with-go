// 10-http-client-basic/main.go
// Cliente HTTP básico en Go
package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Respuesta del servidor:")
	fmt.Println(string(body))
}
