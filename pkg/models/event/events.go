package event

import (
	"database/sql"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/categories"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
)

const (
	tableName = "events"
	schema    = "public"
)

type Event struct {
	EventID     string                   `row:"event_id" type:"exact" json:"event_id"`
	Name        string                   `row:"name" type:"exact" json:"name"`
	Description string                   `row:"description" type:"exact" json:"description"`
	Category    categories.EventCategory `row:"category" type:"exact" fk:"public.event_category" fr:"generated_id"`
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
func (a Model) Create(data Event) error {
	query, args := models.QueryBuilderCreate(data, schema, tableName)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Get data from db into slice of struct
// Searches by the member provided in input struct
func (a Model) Get(data Event) (event []Event) {
	query, args := models.QueryBuilderJoin(data, schema, tableName)

	helpers.LogInfo(args)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &event)
	return
}
