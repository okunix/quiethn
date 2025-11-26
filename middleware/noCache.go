package middleware

import "net/http"

func NoCache() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", "no-cache")
			next.ServeHTTP(w, r)
		})
	}
}
