package routes

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/middleware"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/api"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/authentication"
	"github.com/gorilla/mux"
)

// App is an interface for alll sub routes
type App interface {
	Register(r *mux.Router)
}

var apps = []App{api.API{}, authentication.Authentication{}}

// RegisterApps registers all sub routes
func RegisterApps(r *mux.Router) {
	r.Use(middleware.CorsMiddleware())
	for _, i := range apps {
		i.Register(r)
	}
}
