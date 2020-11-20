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
	StaticPath2  string
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
	configPath = GetExecutableDirectory()
	PgConnString = fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("hostport"), os.Getenv("hostname"), os.Getenv("username"), os.Getenv("password"), os.Getenv("databasename"))
	LogsPath = filepath.Join(GetExecutableDirectory(), "logs")
	StaticPath = filepath.Join(GetExecutableDirectory(), "web/v1")
	StaticPath2 = filepath.Join(GetExecutableDirectory(), "web/v2/build")
	SwaggerPath = filepath.Join(GetExecutableDirectory(), "swagger")
	DatabasePath = filepath.Join(GetExecutableDirectory(), "database")
}
