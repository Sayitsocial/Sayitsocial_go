package middleware

import (
	"net/http"
	"strings"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/authentication"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/common"
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

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return false
	}

	token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
		return helpers.GetJWTKey(), nil
	})
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if token.Valid && token.Method.Alg() == jwt.SigningMethodHS256.Alg() {
		return true
	}
	return false
}
