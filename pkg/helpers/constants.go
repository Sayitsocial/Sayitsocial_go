package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	configPath   string
	LogsPath     string
	StaticPath   string
	DatabasePath string
	PgConnString string
)

const (
	UsernameKey = "username"
	PasswordKey = "password"
	SessionsKey = "sessions"
	PrevURLKey  = "prevurl"
	AuthTypeKey = "type"

	LoginURL   = "/auth/login/"
	HomeURLVol = "/index-vol.html"
	HomeURLOrg = "/index-org.html"

	AuthTypeVol = "vol"
	AuthTypeOrg = "org"

	RowStructTag = "row"
	PKStructTag  = "pk"

	UserAlreadyExistsError  = "User already exists"
	InvalidCredentialsError = "Invalid username or password"
	InvalidUserTypeError    = "Invalid user type"

	HttpSuccessMessage = "success"

	DbSchemaAuth = "auth"
	DbSchemaOrg  = "organisation"

	DbTableAuth             = "auth"
	DbTableOrganisationData = "organisation"
)

func initPaths() {
	configPath = GetWorkingDirectory()
	PgConnString = fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("hostport"), os.Getenv("hostname"), os.Getenv("username"), os.Getenv("password"), os.Getenv("databasename"))
	LogsPath = filepath.Join(GetWorkingDirectory(), "logs")
	StaticPath = filepath.Join(GetWorkingDirectory(), "web/components")
	DatabasePath = filepath.Join(GetWorkingDirectory(), "database")
}
