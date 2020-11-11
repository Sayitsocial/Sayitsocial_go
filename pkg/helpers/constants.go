package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

// Path constants
var (
	configPath   string
	LogsPath     string
	StaticPath   string
	SwaggerPath  string
	DatabasePath string
	PgConnString string
)

// Generic constants
const (
	UsernameKey = "username"
	PasswordKey = "password"
	SessionsKey = "SESSIONID"
	PrevURLKey  = "prevurl"
	AuthTypeKey = "type"
	UserTypeKey = "typeofUser"

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

	HTTPSuccessMessage = "success"

	DbSchemaAuth   = "auth"
	DbSchemaOrg    = "organisation"
	DbSchemaVol    = "volunteer"
	DbSchemaEvents = "events"
	DbSchemaPublic = "public"
)

func initPaths() {
	configPath = GetWorkingDirectory()
	PgConnString = fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("hostport"), os.Getenv("hostname"), os.Getenv("username"), os.Getenv("password"), os.Getenv("databasename"))
	LogsPath = filepath.Join(GetWorkingDirectory(), "logs")
	StaticPath = filepath.Join(GetWorkingDirectory(), "web")
	SwaggerPath = filepath.Join(GetWorkingDirectory(), "swagger")
	DatabasePath = filepath.Join(GetWorkingDirectory(), "database")
}
