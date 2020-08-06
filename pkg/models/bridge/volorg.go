package bridge

import (
	"database/sql"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/orgdata"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/voldata"
)

const (
	tableName = helpers.DbTableVolOrgBridge
	schema    = helpers.DbSchemaBridge
	component = "bridgeModel"
)

type VolOrgRel struct {
	VolunteerID    string `row:"volunteer_id" type:"exact" json:"volunteer_id"`
	OrganisationID string `row:"organisation_id" type:"exact" json:"organisation_id"`
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

func (a Model) Create(data VolOrgRel) error {
	query, args := models.QueryBuilderCreate(data, schema, tableName)

	_, err := a.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (a Model) genericGet(data VolOrgRel) (volOrgRels []VolOrgRel) {
	query, args := models.QueryBuilderGet(data, schema, tableName)
	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}
	models.GetIntoStruct(row, &volOrgRels)
	return
}

// TODO: Do this all in a single query
// This shit is slow af
func (a Model) GetVolunteers(data VolOrgRel) []voldata.VolData {
	fetchedVols := make([]voldata.VolData, 0)

	volOrgRels := a.genericGet(data)

	volModel := voldata.Initialize()
	defer volModel.Close()

	for _, vol := range volOrgRels {
		fetchedVols = append(fetchedVols, volModel.Get(voldata.VolData{VolunteerID: vol.VolunteerID})...)
	}
	return fetchedVols
}

func (a Model) GetOrganisations(data VolOrgRel) []orgdata.OrgData {
	fetchedOrgs := make([]orgdata.OrgData, 0)
	volOrgRels := a.genericGet(data)

	orgModel := orgdata.Initialize()
	defer orgModel.Close()

	for _, vol := range volOrgRels {
		fetchedOrgs = append(fetchedOrgs, orgModel.Get(orgdata.OrgData{OrganisationID: vol.OrganisationID})...)
	}
	return fetchedOrgs
}
