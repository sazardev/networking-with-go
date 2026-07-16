// DNS Lookups Example
// Demonstrates the common net.Lookup* functions: A/AAAA, MX, TXT,
// and a reverse (PTR) lookup.
package main

import (
	"fmt"
	"net"
)

func main() {
	// A/AAAA records: resolve a hostname to its IP addresses.
	ips, err := net.LookupHost("example.com")
	if err != nil {
		fmt.Println("lookup error:", err)
		return
	}
	fmt.Println("A/AAAA records:", ips)

	// MX records: which servers handle mail for this domain.
	mxRecords, err := net.LookupMX("example.com")
	if err == nil {
		for _, mx := range mxRecords {
			fmt.Printf("MX: %s (priority %d)\n", mx.Host, mx.Pref)
		}
	}

	// TXT records: arbitrary text data attached to the domain.
	txtRecords, err := net.LookupTXT("example.com")
	if err == nil {
		fmt.Println("TXT records:", txtRecords)
	}

	// PTR (reverse) lookup: IP address back to hostname.
	names, err := net.LookupAddr("93.184.216.34")
	if err == nil {
		fmt.Println("reverse lookup:", names)
	}
}
