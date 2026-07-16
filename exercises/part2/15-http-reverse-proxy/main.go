// HTTP Reverse Proxy Example
// Forwards every request to a single backend using the standard
// library's production-grade httputil.ReverseProxy.
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	backend, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(backend)

	// Wrap the proxy to add a custom header on every forwarded request,
	// a common pattern for injecting tracing or identifying the proxy hop.
	handler := func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Forwarded-By", "go-reverse-proxy")
		proxy.ServeHTTP(w, r)
	}

	log.Println("reverse proxy listening on :9300, forwarding to :8080")
	log.Fatal(http.ListenAndServe(":9300", http.HandlerFunc(handler)))
}
