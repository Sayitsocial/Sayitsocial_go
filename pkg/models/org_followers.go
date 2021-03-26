package models

import "github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

// swagger:model
type Followers struct {
	GeneratedID    string  `sorm:"generated_id,pk_manual" json:"generated_id"`
	OrganisationID string  `sorm:"organisation_id,pk_manual"  json:"organisation_id"`
	Volunteer      VolData `sorm:"volunteer_id,ft_volunteer.volunteer,fk_volunteer_id" json:"volunteer"`
}

func (Followers) GetTableName() (string, string) {
	return helpers.DbSchemaOrg, "follower_bridge"
}
