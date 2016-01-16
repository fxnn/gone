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

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s %s [%s]", clientIP(r), r.Method, r.URL, time.Since(start))
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
