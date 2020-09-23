package categories

import (
	"database/sql"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
)

const (
	tableName = "event_category"
	schema    = "public"
)

type EventCategory struct {
	GeneratedID string `row:"generated_id" type:"exact" json:"generated_id" pk:"auto"`
	Name        string `row:"name" type:"exact" json:"name"`
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
func (a Model) Create(data EventCategory) error {
	query, args := models.QueryBuilderCreate(data, schema, tableName)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Get data from db into slice of struct
// Searches by the member provided in input structs
func (a Model) Get(data EventCategory) (eventCat []EventCategory) {
	query, args := models.QueryBuilderGet(data, schema, tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &eventCat)
	return
}
