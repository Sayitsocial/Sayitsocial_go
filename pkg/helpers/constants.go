package helpers

import "path/filepath"

var (
	configPath   string
	LogsPath     string
	StaticPath   string
	TemplatePath string
	DatabasePath string
)

const (
	UsernameKey = "username"
	PasswordKey = "password"
	SessionsKey = "sessions"
	LoginURL    = "/auth/login/"
	PrevURLKey  = "prevurl"

	RowStructTag = "row"
	PKStructTag  = "pk"
)

func initPaths() {
	configPath = GetWorkingDirectory()
	LogsPath = filepath.Join(GetWorkingDirectory(), "logs")
	StaticPath = filepath.Join(GetWorkingDirectory(), "web/components/static")
	TemplatePath = filepath.Join(GetWorkingDirectory(), "web/components/html")
	DatabasePath = filepath.Join(GetWorkingDirectory(), "database")
}
