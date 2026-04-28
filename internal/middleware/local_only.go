package middleware

import (
	"net"
	"net/http"
)

// check for IP address
// only let professor go through localhost: 127.0.0.1
// prevent students to poke at professor's endpoint when application is exposed
func LocalOnly(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		if host != "127.0.0.1" && host != "::1" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
