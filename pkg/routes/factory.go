package routes

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/middleware"
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/api"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/authentication"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/sayitsocial"
	"github.com/gorilla/mux"
)

type App interface {
	Register(r *mux.Router)
}

var apps = []App{api.Api{}, authentication.Authentication{}, sayitsocial.Sayitsocial{}}

func RegisterApps(r *mux.Router) {
	r.Use(middleware.RedirectMiddleware())
	for _, i := range apps {
		i.Register(r)
	}
}

func RegisterFileServer(r *mux.Router) {
	r.PathPrefix("/devhtml").Handler(http.StripPrefix("/devhtml",
		http.FileServer(http.Dir(helpers.StaticPath)),
	))

	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/",
		http.FileServer(http.Dir(helpers.SwaggerPath)),
	))
}
