package middleware

import "net/http"

// middleware to tell client browser and Cloudflare's edge cache but check with server first
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

// immutable for files to stay the same for a year without chekcing with server
// request without ?v= falls back to no-cache so un-versioned URL is never pinned
func Immutable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check for ?v= in URL
		// if URL contain "v" , mark it immutable for a year (31536000s)
		if r.URL.Query().Get("v") != "" {
			w.Header().Set("Cache-Control", "public, max-age=31536000, imutable")
		} else {
			// set no-cache for url has no "v"
			w.Header().Set("Cache-Control", "no-cache")
		}
		next.ServeHTTP(w, r)
	})
}
