// DNS Cache Example
// A tiny TTL-respecting DNS cache guarded by a mutex, with negative
// caching for failed lookups.
package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type cacheEntry struct {
	ips     []string
	err     error
	expires time.Time
}

type dnsCache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

func (c *dnsCache) lookup(host string) ([]string, error) {
	c.mu.Lock()
	if entry, ok := c.entries[host]; ok && time.Now().Before(entry.expires) {
		defer c.mu.Unlock()
		if entry.err != nil {
			return nil, entry.err
		}
		return entry.ips, nil
	}
	c.mu.Unlock()

	ips, err := net.LookupHost(host)

	// Negative caching: a lookup failure (NXDOMAIN, a resolver timeout)
	// is cached too, just for a shorter window. Without this, a hostname
	// that doesn't exist gets re-queried on every single call site that
	// asks for it, which is exactly the kind of repeated, avoidable
	// traffic caching was meant to eliminate in the first place.
	ttl := 60 * time.Second
	if err != nil {
		ttl = 5 * time.Second
	}

	c.mu.Lock()
	c.entries[host] = cacheEntry{
		ips:     ips,
		err:     err,
		expires: time.Now().Add(ttl),
	}
	c.mu.Unlock()
	return ips, err
}

func main() {
	cache := &dnsCache{entries: make(map[string]cacheEntry)}

	start := time.Now()
	ips, err := cache.lookup("example.com")
	fmt.Printf("first lookup (%v): ips=%v err=%v\n", time.Since(start), ips, err)

	// Second call within the TTL window should be served from the cache.
	start = time.Now()
	ips, err = cache.lookup("example.com")
	fmt.Printf("cached lookup (%v): ips=%v err=%v\n", time.Since(start), ips, err)
}
