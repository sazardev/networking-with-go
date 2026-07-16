// Rendezvous Server Example
// Matches two UDP peers that register under the same session ID,
// telling each the other's public address so they can hole-punch
// directly — see 17-udp-hole-punching for the matching peer.
package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type pendingPeer struct {
	addr   *net.UDPAddr
	joined time.Time
}

type rendezvous struct {
	mu    sync.Mutex
	peers map[string]pendingPeer // session ID -> first peer's address
}

func (r *rendezvous) handle(conn *net.UDPConn, addr *net.UDPAddr, session string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	other, ok := r.peers[session]
	if !ok {
		// First peer for this session: remember its address and wait.
		r.peers[session] = pendingPeer{addr: addr, joined: time.Now()}
		return
	}

	// Second peer arrived: tell each peer the other's public address.
	// Checking the error matters here — a silently dropped reply means
	// one peer waits forever for an address that was actually sent.
	if _, err := conn.WriteToUDP([]byte(other.addr.String()), addr); err != nil {
		fmt.Println("notify second peer failed:", err)
	}
	if _, err := conn.WriteToUDP([]byte(addr.String()), other.addr); err != nil {
		fmt.Println("notify first peer failed:", err)
	}
	delete(r.peers, session)
}

// sweep evicts sessions whose first peer never got a match within
// maxAge, so an abandoned or mistyped session ID doesn't sit in
// memory for the life of the server.
func (r *rendezvous) sweep(maxAge time.Duration) {
	for range time.Tick(30 * time.Second) {
		r.mu.Lock()
		for session, p := range r.peers {
			if time.Since(p.joined) > maxAge {
				delete(r.peers, session)
			}
		}
		r.mu.Unlock()
	}
}

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 9400})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("rendezvous server listening on :9400")

	r := &rendezvous{peers: make(map[string]pendingPeer)}
	go r.sweep(2 * time.Minute)

	buf := make([]byte, 256)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		session := string(buf[:n])
		r.handle(conn, addr, session)
	}
}
