package initializer

import (
	"flag"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/apps"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/database"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/gorilla/mux"
	"net/http"
)

func Init() error {
	err := initHelpers()
	if err != nil {
		return err
	}

	err = database.RunMigrations()
	if err != nil {
		return err
	}

	err = initWebApp()
	if err != nil {
		return err
	}
	return nil
}

func initHelpers() error {
	err := helpers.ConfigInit()
	if err != nil {
		return err
	}
	err = helpers.CreateDirs()
	if err != nil {
		return err
	}
	helpers.LoggerInit()
	helpers.RndInit()
	return nil
}

func initWebApp() error {
	addr := flag.String("addr", ":8000", "Address of server [default :6969]")
	flag.Parse()

	router := mux.NewRouter()

	apps.RegisterFileServer(router)
	apps.RegisterApps(router)

	helpers.LogInfo("Server starting at "+*addr, "Web")

	err := http.ListenAndServe(*addr, router)
	if err != nil {
		return err
	}

	return nil
}
