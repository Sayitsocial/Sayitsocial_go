package bridge

import (
	"database/sql"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/orgdata"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/voldata"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
)

const (
	tableName = "event_user_bridge"
	schema    = "public"
)

type EventUserBridge struct {
	GeneratedID  string          `row:"generated_id" type:"exact" pk:"auto" json:"generated_id"`
	Organisation orgdata.OrgData `row:"organisation_id" type:"exact" json:"organisation_id" fk:"organisation.organisation" fr:"organisation_id"`
	Volunteer    voldata.VolData `row:"volunteer_id" type:"exact" json:"volunteer_id" fk:"volunteer.volunteer" fr:"volunteer_id"`
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

func (a Model) Create(data EventUserBridge) error {
	query, args := models.QueryBuilderCreate(data, schema, tableName)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (a Model) Get(data EventUserBridge) (eventUserBridge []EventUserBridge) {
	query, args := models.QueryBuilderGet(data, schema, tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &eventUserBridge)
	return
}

func (a Model) GetInner(data EventUserBridge) (eventUserBridge []EventUserBridge) {
	query, args := models.QueryBuilderJoin(data, schema, tableName)
	helpers.LogInfo(query)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}
	models.GetIntoNestedStruct(row, &eventUserBridge)
	return
}
