package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	host := "scanme.nmap.org"
	ports := []int{22, 80, 443, 8080}
	for _, port := range ports {
		address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err != nil {
			fmt.Printf("Port %d closed\n", port)
			continue
		}
		fmt.Printf("Port %d open!\n", port)
		conn.Close()
	}
}
