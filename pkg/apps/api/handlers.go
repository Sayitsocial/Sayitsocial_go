package api

import (
	"encoding/json"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/auth"
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
	model := auth.Initialize()
	defer model.Close()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(model.Get(auth.Auth{}))
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}
