package middleware

import "net/http"

// This middleware is only for use with the / endpoint to catch all undefined endpoints
func NotFound(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
