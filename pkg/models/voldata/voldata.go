package voldata

import (
	"database/sql"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
)

const (
	tableName = helpers.DbTableVolunteerData
	schema    = helpers.DbSchemaVol
)

type VolData struct {
	VolunteerID  string `row:"volunteer_id" type:"exact" pk:"manual" json:"volunteer_id"`
	DisplayName  string `row:"display_name" type:"like" json:"display_name"`
	ContactEmail string `row:"contact_email" type:"like" json:"contact_email"`
	ContactPhone string `row:"contact_phone" type:"like" json:"contact_phone"`
	Bio          string `row:"bio" type:"like" json:"bio"`
}

type Model struct {
	conn *sql.DB
}

// Initialize returns model of db with active connection
func Initialize() *Model {
	return &Model{
		conn: models.GetConn(schema, tableName),
	}
}

// Close closes the connection to db
// Model should not be used after close is called
func (a Model) Close() {
	err := a.conn.Close()
	if err != nil {
		helpers.LogError(err.Error())
	}
}

// Create creates a value in database
func (a Model) Create(data VolData) error {
	query, args := models.QueryBuilderCreate(data, schema, tableName)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Get data from db into slice of struct
// Searches by the member provided in input struct
func (a Model) Get(data VolData) (volData []VolData) {
	query, args := models.QueryBuilderGet(data, schema, tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &volData)
	return
}
