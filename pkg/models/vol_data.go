package models

import "github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

// swagger:model
type VolData struct {
	VolunteerID  string `row:"volunteer_id" type:"exact" pk:"manual" json:"volunteer_id"`
	DisplayName  string `row:"display_name" type:"like" json:"display_name"`
	ContactEmail string `row:"contact_email" type:"like" json:"contact_email"`
	ContactPhone string `row:"contact_phone" type:"like" json:"contact_phone"`
	Bio          string `row:"bio" type:"like" json:"bio"`
	Joined       int64  `row:"joined" type:"exact" json:"joined"`
}

func (VolData) OrgData() (string, string) {
	return helpers.DbSchemaVol, "volunteer"
}
