// Round-Robin Load Balancer Example
// A minimal reverse-proxy load balancer: round-robins across
// healthy backends and rechecks health every 5 seconds.
package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
	"time"
)

type backend struct {
	url     *url.URL
	healthy atomic.Bool
}

type loadBalancer struct {
	backends []*backend
	next     atomic.Uint64
}

func (lb *loadBalancer) pick() *backend {
	// Simple round robin, skipping backends currently marked unhealthy.
	for i := 0; i < len(lb.backends); i++ {
		idx := lb.next.Add(1) % uint64(len(lb.backends))
		b := lb.backends[idx]
		if b.healthy.Load() {
			return b
		}
	}
	return nil
}

func (lb *loadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := lb.pick()
	if b == nil {
		http.Error(w, "no healthy backends", http.StatusServiceUnavailable)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(b.url)
	proxy.ServeHTTP(w, r)
}

func (lb *loadBalancer) healthCheckLoop(client *http.Client) {
	for {
		for _, b := range lb.backends {
			req, err := http.NewRequest(
				http.MethodGet, b.url.String()+"/health", nil,
			)
			if err != nil {
				b.healthy.Store(false)
				continue
			}
			resp, err := client.Do(req)
			ok := err == nil && resp.StatusCode == http.StatusOK
			b.healthy.Store(ok)
			if resp != nil {
				resp.Body.Close()
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: 24-load-balancer <backend-url> [backend-url...]")
		os.Exit(1)
	}

	lb := &loadBalancer{}
	for _, raw := range os.Args[1:] {
		u, err := url.Parse(raw)
		if err != nil {
			fmt.Println("bad backend URL:", raw, err)
			os.Exit(1)
		}
		b := &backend{url: u}
		b.healthy.Store(true)
		lb.backends = append(lb.backends, b)
	}

	go lb.healthCheckLoop(&http.Client{Timeout: 2 * time.Second})

	fmt.Println("load balancer listening on :9600")
	http.ListenAndServe(":9600", lb)
}
