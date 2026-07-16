// UDP Hole Punching Example
// Registers with the rendezvous server under a shared session ID,
// then punches a hole through its own NAT toward the other peer's
// reported public address — run with the session ID as argv[1].
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	session := os.Args[1] // both peers must use the same session ID

	rendezvousAddr, err := net.ResolveUDPAddr(
		"udp", "rendezvous.example.com:9400")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Register with the rendezvous server.
	if _, err := conn.WriteToUDP([]byte(session), rendezvousAddr); err != nil {
		panic(err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	buf := make([]byte, 64)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		panic(fmt.Errorf("no reply from rendezvous server: %w", err))
	}
	peerAddr, err := net.ResolveUDPAddr("udp", string(buf[:n]))
	if err != nil {
		panic(err)
	}
	fmt.Println("peer public address:", peerAddr)

	// connected signals the punch goroutine to stop once we've heard
	// back from the peer at least once — the hole is open, and from
	// here on application traffic itself keeps the NAT mapping alive.
	connected := make(chan struct{})

	// Punch: send directly to the peer's public address. The first few
	// packets may be dropped by the peer's NAT before its own outbound
	// packets have opened a matching hole — that's expected and why
	// both sides keep sending on a short interval.
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-connected:
				return
			case <-ticker.C:
				_, err := conn.WriteToUDP([]byte("ping"), peerAddr)
				if err != nil {
					fmt.Println("punch send failed:", err)
				}
			}
		}
	}()

	announced := false
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, from, err := conn.ReadFromUDP(buf)
		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("shutting down")
				return
			default:
				continue // still waiting for the hole to open
			}
		}
		if !announced {
			close(connected)
			announced = true
		}
		fmt.Printf("received %q from %s\n", buf[:n], from)
	}
}
