package initializer

import (
	"flag"
	"net/http"
	"os"
	"os/exec"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Init initializes the whole app
func Init() error {
	err := initHelpers()
	if err != nil {
		return err
	}

	err = database.RunMigrations()
	if err != nil {
		helpers.LogError(err.Error())
		return err
	}

	err = buildReactApp()
	if err != nil {
		helpers.LogError(err.Error())
		return err
	}

	err = initWebApp()
	if err != nil {
		helpers.LogError(err.Error())
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
	return nil
}

func buildReactApp() error {
	cmd := exec.Command("/usr/bin/python3", "scripts/build_react.py")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	helpers.LogInfo(string(out))
	return nil
}

func initWebApp() error {
	addr := flag.String("addr", "0.0.0.0:8000", "Address of server [default :8000]")
	flag.Parse()

	router := mux.NewRouter()
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	routes.RegisterApps(router)
	routes.RegisterFileServer(router)

	helpers.LogInfo("Server starting at " + *addr)

	err := http.ListenAndServe(*addr, loggedRouter)
	if err != nil {
		return err
	}

	return nil
}
