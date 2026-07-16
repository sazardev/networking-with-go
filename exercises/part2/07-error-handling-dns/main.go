// Error Handling DNS Example
// Resolves a hostname that does not exist and inspects the
// concrete *net.DNSError to tell "not found" apart from a timeout.
package main

import (
	"errors"
	"fmt"
	"net"
)

func main() {
	_, err := net.LookupHost("no-such-hostname.example")
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			fmt.Printf(
				"DNS error for %q (not found=%v, timeout=%v)\n",
				dnsErr.Name, dnsErr.IsNotFound, dnsErr.IsTimeout,
			)
		} else {
			fmt.Printf("Could not resolve host: %v\n", err)
		}
	} else {
		fmt.Println("Host resolved successfully!")
	}
}
