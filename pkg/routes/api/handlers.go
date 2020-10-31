package api

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// API is an empty struct overloading an App interface
type API struct {
}

// Register registers the api router for api endpoints
func (a API) Register(r *mux.Router) {

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.StrictSlash(false)

	apiRouter.Use(middleware.AuthMiddleware())

	apiRouter.HandleFunc("/vol/create", volCreateHandler).Methods("POST")
	apiRouter.HandleFunc("/vol/get", volGetHandler).Methods("GET")

	apiRouter.HandleFunc("/org/create", orgCreateHandler).Methods("POST")
	apiRouter.HandleFunc("/org/get", orgGetHandler).Methods("GET")

	apiRouter.HandleFunc("/event/create", eventCreateHandler).Methods("POST")
	apiRouter.HandleFunc("/event/get", eventGetHandler).Methods("GET")
	apiRouter.HandleFunc("/event/host", eventHostBridge).Methods("GET")
	apiRouter.HandleFunc("/event/attendee", eventAttendeeBridge).Methods("GET")
}
