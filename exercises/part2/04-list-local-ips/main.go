package main

import (
	"fmt"
	"net"
)

func main() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Local IP addresses:")
	for _, addr := range addrs {
		fmt.Println(" -", addr.String())
	}
}
