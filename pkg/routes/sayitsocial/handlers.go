package sayitsocial

import (
	"github.com/gorilla/mux"
)

type Sayitsocial struct {
}

const baseURL = ""

func (s Sayitsocial) Register(r *mux.Router) {
	htmlRouter := r.PathPrefix(baseURL).Subrouter()
	htmlRouter.StrictSlash(true)

	//htmlRouter.Use(middleware.AuthMiddleware())
}
