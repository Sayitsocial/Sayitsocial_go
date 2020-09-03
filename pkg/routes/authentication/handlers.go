package authentication

import (
	"encoding/json"
	"fmt"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var decoder = schema.NewDecoder()

type Authentication struct {
}

type Context struct {
	Error string
}

const (
	baseURL = "/auth"
)

var SessionsStore = sessions.NewCookieStore(helpers.GetSessionsKey(), helpers.GetEncryptionKey())

// register SubRouter
func (a Authentication) Register(r *mux.Router) {
	authRouter := r.PathPrefix(baseURL).Subrouter()

	authRouter.StrictSlash(false)

	authRouter.HandleFunc("/login", loginHandler).Methods("POST")
	authRouter.HandleFunc("/logout", logoutHandler).Methods("POST")
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
/*
 * Handles authenticating user and displaying login page
 * Should redirect to respective page
 * Should update request cookie on successful auth
 * TODO: Implement static login page using go templates
 */
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	session, _ := SessionsStore.Get(r, helpers.SessionsKey)

	// If user is already logged in, don't show login page again until logout
	//if ValidateSession(w, r) {
	//	routes.Redirect(w, r, "/home", routes.StatusFound)
	//	return
	//}

	err := r.ParseForm()
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	var creds LoginReq
	err = decoder.Decode(&creds, r.URL.Query())
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if creds.Username != "" && creds.Password == "" {
		var typeOfUser string
		if userIsValid(creds.Username, creds.Password, &typeOfUser) {

			// Since session key is randomly hashed, its value doesn't matter
			session.Values[helpers.UsernameKey] = creds.Username

			// TODO: Set proper max age
			session.Options.MaxAge = 30 * 60

			err := session.Save(r, w)
			if err != nil {
				helpers.LogError(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Redirect to respective page depending on type of user
			switch typeOfUser {
			case helpers.AuthTypeOrg:
				http.Redirect(w, r, helpers.HomeURLOrg, http.StatusOK)
				return
			case helpers.AuthTypeVol:
				http.Redirect(w, r, helpers.HomeURLVol, http.StatusOK)
				return
			}
			return
		}

		// TODO: Should display error on page
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		err := encoder.Encode(&invalidCredentialsError{
			Body: Body{Message: helpers.InvalidCredentialsError},
		})
		if err != nil {
			helpers.LogError(err.Error())
		}
		return
	}
}

// swagger:route POST /auth/logout auth logout
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
	_, err = fmt.Fprintf(w, helpers.HttpSuccessMessage)
	if err != nil {
		helpers.LogError(err.Error())
	}
}

// Validate user from hashes password
func userIsValid(username string, password string, typeOfUser *string) bool {
	model := auth.Initialize()
	defer model.Close()

	fetchUsers := model.Get(auth.Auth{Username: username})
	if len(fetchUsers) > 0 {
		hashedPass := fetchUsers[0].Password
		if hashedPass != "" {
			err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
			if err != nil {
				helpers.LogError(err.Error())
				return false
			}
			*typeOfUser = fetchUsers[0].TypeOfUser
			return true
		}
	}
	return false
}

/* Compares session key to server and validates
 * If valid, then rebuilds session
 */
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
		model := auth.Initialize()
		defer model.Close()

		user := model.Get(auth.Auth{Username: val.(string)})
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

//func IsAdminFromSession(r *routes.Request) bool {
//	session, err := SessionsStore.Get(r, helpers.SessionsKey)
//	if err != nil {
//		return false
//	}
//	if session.IsNew {
//		return false
//	}
//	model := authvol.Initialize()
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
