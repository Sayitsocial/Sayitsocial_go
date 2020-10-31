package middleware

import (
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/authentication"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/common"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func AuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if CookieAuthCheck(w, r) || JWTCheck(w, r) {
				next.ServeHTTP(w, r)
			} else {
				common.WriteError("Unauthorized", http.StatusUnauthorized, w)
			}
		})
	}
}

// CookieAuthCheck authenticates user from cookies
func CookieAuthCheck(w http.ResponseWriter, r *http.Request) bool {
	return authentication.ValidateSession(w, r)
}

// JWTCheck authenticates user from JW token
func JWTCheck(w http.ResponseWriter, r *http.Request) bool {
	if len(helpers.GetJWTKey()) == 0 {
		helpers.LogError("JWT Auth key empty")
		return false
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return (helpers.GetJWTKey()), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	if jwtMiddleware.CheckJWT(w, r) != nil {
		return true
	}
	return false
}
