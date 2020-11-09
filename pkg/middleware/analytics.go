package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// AnalyticsMiddleware is responsible for tracking user behaviour through request urls
func AnalyticsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}
}
