package main

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers/initializer"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	helpers.SetWorkingDirectory(wd)
	err = initializer.Init()
	if err != nil {
		panic(err)
	}
}
