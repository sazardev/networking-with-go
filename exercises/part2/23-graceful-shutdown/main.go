// Graceful Shutdown Example
// Combines signal.NotifyContext with http.Server.Shutdown so
// in-flight requests finish before the process exits on SIGINT/SIGTERM.
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func buildHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello\n"))
	})
	return mux
}

func main() {
	srv := &http.Server{Addr: ":8080", Handler: buildHandler()}

	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	log.Println("server started on :8080")

	<-ctx.Done() // blocks until SIGINT or SIGTERM arrives
	log.Println("shutdown signal received, draining connections...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(), 15*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
	log.Println("server stopped cleanly")
}
