package helpers

import (
	"os"
	"path/filepath"
	"strconv"
)

var basePath = ""

// GetWorkingDirectory returns the set basePath
func GetWorkingDirectory() string {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if !debug || err != nil {
		return GetExecutableDirectory()
	}
	return basePath

}

// GetExecutableDirectory returns the directory of executable
func GetExecutableDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		LogError(err.Error())
	}
	return dir
}

// SetWorkingDirectory sets basePath to provided value
func SetWorkingDirectory(path string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		LogError(err.Error())
		abs = "."
	}
	basePath = abs
}
