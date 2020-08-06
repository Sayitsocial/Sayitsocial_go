package orgdata

import (
	"database/sql"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
)

const (
	tableName = helpers.DbTableOrganisationData
	schema    = helpers.DbSchemaOrg
	component = "orgModel"
)

type OrgData struct {
	OrganisationID string `row:"organisation_id" type:"exact" pk:"manual" json:"organisation_id"`
	DisplayName    string `row:"display_name" type:"like" json:"display_name"`
	Locality       string `row:"locality" type:"like" json:"locality"`
	RegistrationNo string `row:"registration_no" type:"exact" json:"registration_no"`
	ContactEmail   string `row:"contact_email" type:"like" json:"contact_email"`
	ContactPhone   string `row:"contact_phone" type:"like" json:"contact_phone"`
	Desc           string `row:"description" type:"like" json:"desc"`
	Owner          string `row:"owner" type:"like" json:"owner"`
	Achievements   string `row:"achievements" type:"like" json:"achievements"`
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
		helpers.LogError(err.Error(), component)
	}
}

func (a Model) Create(data OrgData) error {
	query, args := models.QueryBuilderCreate(data, schema, tableName)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (a Model) Get(data OrgData) (orgData []OrgData) {
	query, args := models.QueryBuilderGet(data, schema, tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}

	models.GetIntoStruct(row, &orgData)
	return
}