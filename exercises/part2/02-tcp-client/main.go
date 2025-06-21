package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	buf := make([]byte, 4096)
	n, _ := conn.Read(buf)
	fmt.Println(string(buf[:n]))
	conn.Close()
}
