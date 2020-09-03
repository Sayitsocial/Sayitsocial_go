package middleware

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/gorilla/mux"
	"net/http"
)

func RedirectMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			helpers.LogInfo(r.URL.Path)
			if r.URL.Path == "/" {
				helpers.LogInfo("redirecting")
				http.Redirect(w, r, helpers.LoginURL, http.StatusFound)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
