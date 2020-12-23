package models

import (
	"encoding/json"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder/types"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

// OrgData is a model to store data about an organisation
// swagger:model
type OrgData struct {
	OrganisationID string                `row:"organisation_id" type:"exact" pk:"manual" json:"organisation_id"`
	DisplayName    string                `row:"display_name" type:"like" json:"display_name"`
	RegistrationNo string                `row:"registration_no" type:"exact" json:"registration_no,omitempty"`
	ContactEmail   string                `row:"contact_email" type:"like" json:"contact_email,omitempty"`
	ContactPhone   string                `row:"contact_phone" type:"like" json:"contact_phone"`
	Desc           string                `row:"description" type:"like" json:"desc,omitempty"`
	Owner          string                `row:"owner" type:"like" json:"owner,omitempty"`
	Achievements   string                `row:"achievements" type:"like" json:"achievements,omitempty"`
	TypeOfOrg      int                   `row:"type_of_org" type:"like" json:"type_of_org"`
	Followers      uint64                `row:"followers" type:"exact" json:"follower_count"`
	Location       types.GeographyPoints `row:"location" type:"onlyvalue" json:"location"`
	Short          bool                  `scan:"ignore" json:"-"`
}

func (OrgData) OrgData() (string, string) {
	return helpers.DbSchemaOrg, "organisation"
}

// MarshalJSON is responsible for custom marshaling struct to json
func (o *OrgData) MarshalJSON() ([]byte, error) {
	type tmp OrgData
	//cat := &e.Category
	helpers.LogInfo(o.Short)
	if o.Short {
		o.RegistrationNo = ""
		o.ContactEmail = ""
		o.ContactPhone = ""
		o.Desc = ""
		o.Owner = ""
		o.Achievements = ""
	}
	return json.Marshal(&struct {
		*tmp
	}{
		(*tmp)(o),
	})
}
