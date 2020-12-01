package categories

import (
	"database/sql"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

const (
	tableName = "event_category"
	schema    = helpers.DbSchemaEvents
)

type EventCategory struct {
	GeneratedID int    `row:"generated_id" type:"exact" json:"generated_id,omitempty" pk:"auto"`
	Name        string `row:"name" type:"exact" json:"name,omitempty"`
}

type Model struct {
	trans *sql.Tx
	conn  *sql.DB
}

// Initialize returns model of db with active connection
func Initialize(tx *sql.Tx) *Model {
	if tx != nil {
		return &Model{
			trans: tx,
		}
	}
	return &Model{
		conn: querybuilder.GetConn(schema, tableName),
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
	query, args := querybuilder.QueryBuilderCreate(data, schema, tableName)

	var err error
	if a.trans != nil {
		_, err = a.trans.Exec(query, args...)
	} else {
		_, err = a.conn.Exec(query, args...)
	}
	return err
}

// Get data from db into slice of struct
// Searches by the member provided in input structs
func (a Model) Get(data EventCategory) (eventCat []EventCategory) {
	query, args := querybuilder.QueryBuilderGet(data, schema+"."+tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	querybuilder.GetIntoStruct(row, &eventCat)
	return
}
