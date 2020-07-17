package sayitsocial

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
)

type Jizzberry struct {
}

const component = "Web"
const baseURL = ""

func (a Jizzberry) Register(r *mux.Router) {
	htmlRouter := r.PathPrefix(baseURL).Subrouter()
	htmlRouter.StrictSlash(true)

	//htmlRouter.Use(middleware.AuthMiddleware())

	htmlRouter.HandleFunc("/login-org", loginOrgHandler)
	htmlRouter.HandleFunc("/login-volunteer", loginVolHandler)
	htmlRouter.HandleFunc("/signup-volunteer", registerVolHandler)
	htmlRouter.HandleFunc("/choose-identity", chooseIdentityHandler)
}

func loginOrgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, filepath.Join(helpers.TemplatePath, "login-org.html"))
}

func loginVolHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, filepath.Join(helpers.TemplatePath, "login-volunteer.html"))
}

func registerVolHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, filepath.Join(helpers.TemplatePath, "signup-volunteer.html"))
}

func chooseIdentityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, filepath.Join(helpers.TemplatePath, "choose-identity.html"))
}
