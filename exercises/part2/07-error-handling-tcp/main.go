// Error Handling TCP Example
// Tries to connect to a TCP server and handles connection errors.
package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9999") // Try to connect to a port with no server
	if err != nil {
		fmt.Println("Error connecting:", err) // Print the error
		return
	}
	fmt.Fprintln(conn, "Hello!")
	conn.Close()
}
