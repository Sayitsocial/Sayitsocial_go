package sayitsocial

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type Jizzberry struct {
}

const component = "Web"
const baseURL = "/Jizzberry"

func (a Jizzberry) Register(r *mux.Router) {
	htmlRouter := r.PathPrefix(baseURL).Subrouter()
	htmlRouter.StrictSlash(true)

	htmlRouter.Use(middleware.AuthMiddleware())

	htmlRouter.HandleFunc("/home", homeHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
}
