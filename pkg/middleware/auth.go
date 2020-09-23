package middleware

import (
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/authentication"
	"github.com/gorilla/mux"
)

// CookieAuthMiddleware authenticates user from cookies
func CookieAuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authentication.ValidateSession(w, r) {
				next.ServeHTTP(w, r)
			} else {
				session, _ := authentication.SessionsStore.Get(r, helpers.SessionsKey)
				session.Values[helpers.PrevURLKey] = r.URL.Path
				err := session.Save(r, w)
				if err != nil {
					helpers.LogError(err.Error())
				}
				http.Redirect(w, r, helpers.LoginURL, http.StatusFound)
			}
		})
	}
}
