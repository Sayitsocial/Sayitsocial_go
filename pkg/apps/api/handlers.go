package api

import (
	"fmt"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/gorilla/mux"
	"net/http"
)

const component = "API"

type Api struct {
}

func (a Api) Register(r *mux.Router) {

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.StrictSlash(false)

	//apiRouter.Use(middleware.AuthMiddleware())

	apiRouter.HandleFunc("/dummy", dummyHandler).Methods("GET")
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Success")
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}
