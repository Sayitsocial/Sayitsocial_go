package models

import "github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

// swagger:model
type Followers struct {
	GeneratedID    string  `row:"generated_id" type:"exact" pk:"manual" json:"generated_id"`
	OrganisationID string  `row:"organisation_id" type:"exact" pk:"manual" json:"organisation_id"`
	Volunteer      VolData `row:"volunteer_id" type:"exact" ft:"volunteer.volunteer" fk:"volunteer_id" json:"volunteer"`
}

func (Followers) GetTableName() (string, string) {
	return helpers.DbSchemaOrg, "follower_bridge"
}
