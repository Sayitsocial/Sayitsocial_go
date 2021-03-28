package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var decoder = schema.NewDecoder()

// Authentication is an empty struct overloading app interface
type Authentication struct {
}

const (
	baseURL = "/auth"
)

// SessionsStore memory store for cookies
var SessionsStore = sessions.NewCookieStore(helpers.GetSessionsKey(), helpers.GetEncryptionKey())

// Register registers SubRouter
func (a Authentication) Register(r *mux.Router) {
	authRouter := r.PathPrefix(baseURL).Subrouter()

	authRouter.StrictSlash(false)

	authRouter.HandleFunc("/login", loginHandler).Methods("POST", "GET")
	authRouter.HandleFunc("/logout", logoutHandler).Methods("POST")
	authRouter.HandleFunc("/jwt-login", jwtLoginHandler).Methods("POST")
	authRouter.HandleFunc("/jwt-refresh", jwtRefreshHandler).Methods("POST")
	authRouter.HandleFunc("/isLogged", isLogged).Methods("GET")
}

// swagger:route POST /auth/login auth login
//
// Login to existing account
//
//
//     Consumes:
//     - application/json
//
//
//     Schemes: http
//
//
//     Security:
//
//     Responses:
//       200: successResponse
//		 401: unauthorizedError
/*
 * Handles authenticating user and displaying login page
 * Should redirect to respective page
 * Should update request cookie on successful auth
 * TODO: Implement static login page using go templates
 */
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	session, _ := SessionsStore.Get(r, helpers.SessionsKey)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var creds LoginReq
	json.Unmarshal(body, &creds)
	helpers.LogInfo(creds)
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if creds.Username != "" && creds.Password != "" {
		if ok, typeOfUser := isUserValid(creds.Username, creds.Password); ok {

			// Since session key is randomly hashed, its value doesn't matter
			session.Values[helpers.UsernameKey] = creds.Username
			session.Values[helpers.UserTypeKey] = typeOfUser

			// TODO: Set proper max age
			session.Options.MaxAge = int((30 * time.Minute).Seconds())

			err := session.Save(r, w)
			if err != nil {
				helpers.LogError(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			common.WriteSuccess(w)
			return
		}

		// TODO: Should display error on page
		err := common.WriteError("invalid credentials", http.StatusUnauthorized, w)
		if err != nil {
			helpers.LogError(err.Error())
		}
		return
	}
}

// jwtLoginHandler
// swagger:route POSt /auth/jwt-login auth JWTLogin
//
// Login to existing account
//
//
//     Consumes:
//     - application/json
//
//
//     Schemes: http
//
//
//     Security:
//
//     Responses:
//       200: JWTLoginResp
//		 401: unauthorizedError
func jwtLoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds LoginReq
	err := decoder.Decode(&creds, r.URL.Query())
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.LogInfo(creds)
	if creds.Username != "" && creds.Password != "" {
		if ok, _ := isUserValid(creds.Username, creds.Password); ok {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				helpers.UsernameKey: creds.Username,
				"expiry":            time.Now().Add(time.Hour * time.Duration(1)).Unix(),
				"iat":               time.Now().Unix(),
			})
			tokenString, err := token.SignedString(helpers.GetJWTKey())
			if err != nil {
				err := common.WriteError("Token generation failed", http.StatusInternalServerError, w)
				if err != nil {
					helpers.LogError(err.Error())
				}
				return
			}
			err = json.NewEncoder(w).Encode(JWTResp{
				Token: tokenString,
			})

			helpers.LogInfo(tokenString)
			if err != nil {
				helpers.LogError(err.Error())
				common.WriteError(err.Error(), http.StatusInternalServerError, w)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	err = common.WriteError("invalid credentials", http.StatusUnauthorized, w)
	if err != nil {
		helpers.LogError(err.Error())
	}
	return
}

// jwtRefreshHandler
// swagger:route POST /auth/jwt-refresh auth JWTRefresh
//
// Login to existing account
//
//
//     Consumes:
//     - application/json
//
//
//     Schemes: http
//
//
//     Security:
//
//     Responses:
//       200: JWTLoginResp
//		 401: unauthorizedError
func jwtRefreshHandler(w http.ResponseWriter, r *http.Request) {
	var tokenPrev JWTRefreshReq
	err := decoder.Decode(&tokenPrev, r.URL.Query())
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := jwt.ParseWithClaims(tokenPrev.Token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() == jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Method.Alg())
		}
		return helpers.GetJWTKey(), nil
	})
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			helpers.UsernameKey: claims[helpers.UsernameKey],
			"expiry":            time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			"iat":               time.Now().Unix(),
		})
		tokenString, err := newToken.SignedString(helpers.GetJWTKey())
		if err != nil {
			err := common.WriteError("Token generation failed", http.StatusInternalServerError, w)
			if err != nil {
				helpers.LogError(err.Error())
			}
			return
		}
		err = json.NewEncoder(w).Encode(JWTResp{
			Token: tokenString,
		})

		if err != nil {
			helpers.LogError(err.Error())
			common.WriteError(err.Error(), http.StatusInternalServerError, w)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	}
}

// swagger:route POST /auth/logout auth logout
//
// Logout from existing account
//
//
//     Consumes:
//     - application/json
//
//
//     Schemes: http
//
//
//     Security:
//
//     Responses:
//       200: successResponse
// Deletes session on logout
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, helpers.UsernameKey)

	// Instantly expires session
	session.Options.MaxAge = -1

	err = session.Save(r, w)

	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	common.WriteSuccess(w)
}

// Validate user from hashes password
func isUserValid(username string, password string) (bool, string) {
	model, err := querybuilder.Initialize(nil, nil)
	if err != nil {
		helpers.LogError(err.Error())
	}
	defer model.Close()

	x, err := model.Get(models.Auth{
		Username: username,
	})
	if err != nil {
		helpers.LogError(err.Error())
		return false, ""
	}

	fetchUsers := *x.(*[]models.Auth)
	if len(fetchUsers) > 0 {
		hashedPass := fetchUsers[0].Password
		if hashedPass != "" {
			err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
			if err != nil {
				helpers.LogError(err.Error())
				return false, ""
			}
			return true, fetchUsers[0].TypeOfUser
		}
	}
	return false, ""
}

// ValidateSession compares session key to server and validates
// If valid, then rebuilds session
func ValidateSession(w http.ResponseWriter, r *http.Request) bool {
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		return false
	}

	if session.IsNew {
		return false
	}

	val := session.Values[helpers.UsernameKey]

	if val != nil {
		model, err := querybuilder.Initialize(nil, nil)
		if err != nil {
			helpers.LogError(err.Error())
		}
		defer model.Close()

		x, err := model.Get(models.Auth{Username: val.(string)})
		if err != nil {
			helpers.LogError(err.Error())
			return false
		}
		user := *x.(*[]models.Auth)
		if len(user) > 0 {
			if user[0].Username == val {
				session.Options.MaxAge = 30 * 60
				err := session.Save(r, w)
				if err != nil {
					helpers.LogError(err.Error())
					return false
				}
				return true
			}
		}
	}
	return false
}

func isLogged(w http.ResponseWriter, r *http.Request) {
	ok := ValidateSession(w, r)
	_, err := fmt.Fprintf(w, "%v", ok)
	if err != nil {
		helpers.LogError(err.Error())
	}
}

// GetUsernameFromSession returns username from http request
func GetUsernameFromSession(r *http.Request) string {
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		return ""
	}
	if session.IsNew {
		return ""
	}

	return session.Values[helpers.UsernameKey].(string)
}

// func IsAdminFromSession(r *routes.Request) bool {
//	session, err := SessionsStore.Get(r, helpers.SessionsKey)
//	if err != nil {
//		return false
//	}
//	if session.IsNew {
//		return false
//	}
//	model := querybuilder.Initialize()
//	defer model.Close()
//
//	user := model.Get(authvol.AuthVol{Username: session.Values[helpers.UsernameKey].(string)})
//	if len(user) > 0 {
//		if user[0].IsAdmin {
//			return true
//		}
//	}
//	return false
//}
