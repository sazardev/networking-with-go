// UDP Broadcast Example
// Sends a broadcast message to all devices on the local network.
package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. net.Dial with broadcast address sends to all devices
	conn, err := net.Dial("udp", "255.255.255.255:9002")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Fprintf(conn, "UDP broadcast message!\n")
	conn.Close()
}
