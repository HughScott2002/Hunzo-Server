package middleware

import (
	"net/http"

	"example.com/m/v2/src/utils"
)

// RateLimitMiddleware is a middleware handler that performs rate limiting
func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		limiter := utils.GetVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}
