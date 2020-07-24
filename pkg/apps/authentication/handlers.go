package authentication

import (
	"fmt"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Authentication struct {
}

type Context struct {
	Error string
}

const (
	baseURL   = "/auth"
	component = "WebAuth"
)

var SessionsStore = sessions.NewCookieStore(helpers.GetSessionsKey())

func (a Authentication) Register(r *mux.Router) {
	authRouter := r.PathPrefix(baseURL).Subrouter()

	authRouter.StrictSlash(false)

	authRouter.HandleFunc("/login", loginHandler).Methods("POST")
	authRouter.HandleFunc("/logout", logoutHandler).Methods("POST")
	authRouter.HandleFunc("/create", newUser).Methods("POST")
	authRouter.HandleFunc("/isLogged", isLogged).Methods("GET")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	session, _ := SessionsStore.Get(r, helpers.SessionsKey)

	// If user is already logged in, don't show login page again until logout
	//if ValidateSession(w, r) {
	//	http.Redirect(w, r, "/home", http.StatusFound)
	//	return
	//}

	queryParams := r.URL.Query()

	err := r.ParseForm()
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	if user, ok := queryParams[helpers.UsernameKey]; ok && len(user) > 0 {
		if pass, ok := queryParams[helpers.PasswordKey]; ok && len(pass) > 0 {
			username := user[0]
			password := pass[0]

			if username != "" && password != "" {
				if userIsValid(username, password) {
					session.Values[helpers.UsernameKey] = username
					session.Options.MaxAge = 30 * 60

					err := session.Save(r, w)
					if err != nil {
						helpers.LogError(err.Error(), component)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					_, err = fmt.Fprintf(w, "Success")
					if err != nil {
						helpers.LogError(err.Error(), component)
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				}
				_, err := fmt.Fprintf(w, "Invalid username or password")
				if err != nil {
					helpers.LogError(err.Error(), component)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		helpers.LogError(err.Error(), component)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, helpers.UsernameKey)
	session.Options.MaxAge = -1

	err = session.Save(r, w)

	if err != nil {
		helpers.LogError(err.Error(), component)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintf(w, "success")
	if err != nil {
		helpers.LogError(err.Error(), component)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userIsValid(username string, password string) bool {
	model := auth.Initialize()
	defer model.Close()

	fetchUsers := model.Get(auth.Auth{Username: username})
	if len(fetchUsers) > 0 {
		hashedPass := fetchUsers[0].Password
		if hashedPass != "" {
			err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
			if err != nil {
				helpers.LogError(err.Error(), component)
				return false
			}
			return true
		}
	}
	return false
}

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
					helpers.LogError(err.Error(), component)
					return false
				}
				return true
			}
		}
	}
	return false
}

func newUser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if userval, ok := queryParams[helpers.UsernameKey]; ok && len(userval) > 0 {
		if passval, ok := queryParams[helpers.PasswordKey]; ok && len(passval) > 0 {
			username := userval[0]
			password := passval[0]

			model := auth.Initialize()
			defer model.Close()

			if val := model.Get(auth.Auth{Username: username}); len(val) > 0 {
				http.Error(w, helpers.UserAlreadyExistsError, http.StatusInternalServerError)
				return
			}

			err := model.Create(auth.Auth{
				Username: username,
				Password: password,
				IsAdmin:  false,
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func isLogged(w http.ResponseWriter, r *http.Request) {
	ok := ValidateSession(w, r)
	_, err := fmt.Fprintf(w, "%v", ok)
	if err != nil {
		helpers.LogError(err.Error(), component)
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

func IsAdminFromSession(r *http.Request) bool {
	session, err := SessionsStore.Get(r, helpers.SessionsKey)
	if err != nil {
		return false
	}
	if session.IsNew {
		return false
	}
	model := auth.Initialize()
	defer model.Close()

	user := model.Get(auth.Auth{Username: session.Values[helpers.UsernameKey].(string)})
	if len(user) > 0 {
		if user[0].IsAdmin {
			return true
		}
	}
	return false
}
