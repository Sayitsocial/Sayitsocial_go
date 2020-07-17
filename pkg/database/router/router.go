package router

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"io/ioutil"
	"path/filepath"
)

func GetDatabase(table string) string {
	switch table {

	case "auth":
		return filepath.FromSlash(helpers.DatabasePath + "/auth.db")
	}

	file, _ := ioutil.TempFile(helpers.DatabasePath, "/tmp.db")
	return file.Name()
}
