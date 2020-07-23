package sayitsocial

import (
	"github.com/gorilla/mux"
)

type Jizzberry struct {
}

const component = "Web"
const baseURL = ""

func (a Jizzberry) Register(r *mux.Router) {
	htmlRouter := r.PathPrefix(baseURL).Subrouter()
	htmlRouter.StrictSlash(true)

	//htmlRouter.Use(middleware.AuthMiddleware())
}
