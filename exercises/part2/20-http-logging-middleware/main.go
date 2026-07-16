// HTTP Logging Middleware Example
// Wraps every request with structured logging, capturing the
// response status code via a small statusRecorder wrapper.
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

// statusRecorder wraps http.ResponseWriter to capture the status code,
// since the standard interface has no way to read it back afterward.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		slog.Info("request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello\n"))
	})

	slog.Info("server starting", "addr", ":8083")
	http.ListenAndServe(":8083", loggingMiddleware(mux))
}
