package main

import (
	"fmt"
	"net"
)

func main() {
    addr := "192.168.1.10:8080"
    host, port, err := net.SplitHostPort(addr)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Host:", host)
    fmt.Println("Port:", port)
}
