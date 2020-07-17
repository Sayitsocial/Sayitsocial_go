package helpers

import (
	"os"
	"path/filepath"
)

func CreateDirs() error {
	err := makeDir(filepath.FromSlash(DatabasePath))
	if err != nil {
		return err
	}
	err = makeDir(filepath.FromSlash(LogsPath))
	if err != nil {
		return err
	}
	return nil
}

func makeDir(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
