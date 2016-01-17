package http

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var start = time.Now()
		var stats = &responseWriterWithStats{wrapped: w}

		next.ServeHTTP(stats, r)

		log.Printf("[%s] %d %s %s [%d bytes in %s]", clientIP(r), stats.status, r.Method, r.URL, stats.bytesWritten, time.Since(start))
	})
}

func clientIP(r *http.Request) string {
	if r.Header.Get("X-Forwarded-For") != "" {
		// Proxied requests
		var values = r.Header[http.CanonicalHeaderKey("X-Forwarded-For")]
		return strings.Join(values, ", ")
	}
	return r.RemoteAddr
}

type responseWriterWithStats struct {
	wrapped      http.ResponseWriter
	status       int
	bytesWritten int
}

func (w *responseWriterWithStats) Header() http.Header {
	return w.wrapped.Header()
}

func (w *responseWriterWithStats) Write(content []byte) (written int, err error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	written, err = w.wrapped.Write(content)
	w.bytesWritten += written
	return
}

func (w *responseWriterWithStats) WriteHeader(status int) {
	w.wrapped.WriteHeader(status)
	w.status = status
}
