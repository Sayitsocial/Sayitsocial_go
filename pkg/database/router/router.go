package router

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"io/ioutil"
	"path/filepath"
)

func GetDatabase(table string) string {
	var databaseDir = helpers.GetWorkingDirectory() + "/assets/database"
	switch table {

	case "auth":
		return filepath.FromSlash(databaseDir + "/auth.db")
	}

	file, _ := ioutil.TempFile(databaseDir, "/tmp.db")
	return file.Name()
}
