package main

import (
	"fmt"
	"net"
)

func main() {
	addr := "[2001:db8::1]:443"
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("IPv6 Host:", host)
	fmt.Println("Port:", port)
	ip := net.ParseIP(host)
	if ip != nil && ip.To16() != nil {
		fmt.Println("It's a valid IPv6 address!")
	}
}
