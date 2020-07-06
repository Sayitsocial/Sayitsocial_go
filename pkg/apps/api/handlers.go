package api

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

const component = "API"

type Api struct {
}

type task struct {
	Uid string `json:"uid"`
}

func (a Api) Register(r *mux.Router) {

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.StrictSlash(false)

	apiRouter.Use(middleware.AuthMiddleware())

	apiRouter.HandleFunc("/dummy", dummyHandler).Methods("GET")
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {

}
