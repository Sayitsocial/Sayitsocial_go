package auth

import (
	"database/sql"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	tableName = helpers.DbTableAuth
	schema    = helpers.DbSchemaAuth
)

type Auth struct {
	UID        string `row:"uid" type:"exact"`
	Username   string `row:"username" type:"exact"`
	Password   string `row:"password" type:"exact"`
	TypeOfUser string `row:"typeOfUser" type:"exact"`
}

type Model struct {
	conn *sql.DB
}

func Initialize() *Model {
	return &Model{
		conn: models.GetConn(schema, tableName),
	}
}

func (a Model) Close() {
	err := a.conn.Close()
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (a Model) Create(auth Auth) error {
	auth.Password = hashPassword(auth.Password)

	query, args := models.QueryBuilderCreate(auth, schema, tableName)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (a Model) Get(auth Auth) (allAuth []Auth) {
	query, args := models.QueryBuilderGet(auth, schema, tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &allAuth)
	return
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		helpers.LogError(err.Error())
	}

	return string(hash)
}
