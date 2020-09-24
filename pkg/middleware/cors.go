package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// CorsMiddleware adds cors headers to response
func CorsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			addCorsHeader(w)
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}
