package main

import (
	"fmt"
	"net"
)

func main() {
	// Get all IP addresses for a hostname
	host := "google.com"
	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("IP addresses for", host, ":")
	for _, ip := range ips {
		fmt.Println(" -", ip)
	}
}
