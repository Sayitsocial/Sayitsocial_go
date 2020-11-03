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

// swagger:model
type Event struct {
	EventID     string                   `row:"event_id" type:"exact" json:"event_id" pk:"manual"`
	Name        string                   `row:"name" type:"like" json:"name"`
	Description string                   `row:"description" type:"like" json:"description"`
	StartTime   int64                    `row:"start_time" type:"exact" json:"start_time"`
	HostTime    int64                    `row:"host_time" type:"exact" json:"host_time"`
	Location    models.GeographyPoints   `row:"location" type:"exact" json:"location"`
	Category    categories.EventCategory `row:"category" type:"exact" fk:"public.event_category" fr:"generated_id" json:"category"`
}

type Model struct {
	conn  *sql.DB
	trans *sql.Tx
}

// Initialize returns model of db with active connection
func Initialize(tx *sql.Tx) *Model {
	if tx != nil {
		return &Model{
			trans: tx,
		}
	}
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
	query, args := models.QueryBuilderCreate(data, schema+"."+tableName)
	helpers.LogInfo(args)
	var err error
	if a.trans != nil {
		_, err = a.trans.Exec(query, args...)
	} else {
		_, err = a.conn.Exec(query, args...)
	}
	return err
}

// Get data from db into slice of struct
// Searches by the member provided in input struct
func (a Model) Get(data Event) (event []Event) {
	query, args := models.QueryBuilderJoin(data, schema+"."+tableName)
	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &event)
	return
}

// Count gets count of rows corresponsing to provided search params
func (a Model) Count(data Event) (count []int) {
	query, args := models.QueryBuilderCount(data, schema+"."+tableName)
	helpers.LogInfo(query)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoVar(row, &count)
	return
}
