package sayitsocial

import (
	"github.com/gorilla/mux"
)

type Sayitsocial struct {
}

const baseURL = ""

// Register registers the html router for static webpages
func (s Sayitsocial) Register(r *mux.Router) {
	htmlRouter := r.PathPrefix(baseURL).Subrouter()
	htmlRouter.StrictSlash(true)

	// htmlRouter.Use(middleware.AuthMiddleware())

}
