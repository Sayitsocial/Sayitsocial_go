package apps

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/apps/api"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/apps/authentication"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/apps/sayitsocial"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/gorilla/mux"
	"net/http"
)

type App interface {
	Register(r *mux.Router)
}

var apps = []App{api.Api{}, authentication.Authentication{}, sayitsocial.Jizzberry{}}

func RegisterApps(r *mux.Router) {
	for _, i := range apps {
		i.Register(r)
	}
}

func RegisterFileServer(r *mux.Router) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(helpers.StaticPath)),
	))
}
