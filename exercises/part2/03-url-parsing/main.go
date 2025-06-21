package main

import (
	"fmt"
	"net/url"
)

// Example of parsing a URL and extracting its components.
func main() {
	rawURL := "https://example.com:8080/path?query=go#section"
	parsed, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	fmt.Println("Scheme:", parsed.Scheme)
	fmt.Println("Host:", parsed.Host)
	fmt.Println("Path:", parsed.Path)
	fmt.Println("Query:", parsed.RawQuery)
	fmt.Println("Fragment:", parsed.Fragment)
}
