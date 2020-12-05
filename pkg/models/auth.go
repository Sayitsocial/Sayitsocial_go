package models

import "github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

type Auth struct {
	UID        string `row:"uid" type:"exact"`
	Username   string `row:"username" type:"exact"`
	Password   string `row:"password" type:"exact"`
	TypeOfUser string `row:"typeOfUser" type:"exact"`
}

func (Auth) GetTableName() (string, string) {
	return helpers.DbSchemaAuth, "auth"
}
