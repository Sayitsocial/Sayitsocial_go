package models

import (
	"encoding/json"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder/types"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

// OrgData is a model to store data about an organisation
// swagger:model
type OrgData struct {
	OrganisationID string                `sorm:"organisation_id,pk_manual" json:"organisation_id"`
	DisplayName    string                `sorm:"display_name" json:"display_name"`
	RegistrationNo string                `sorm:"registration_no" json:"registration_no,omitempty"`
	ContactEmail   string                `sorm:"contact_email" json:"contact_email,omitempty"`
	ContactPhone   string                `sorm:"contact_phone" json:"contact_phone"`
	Desc           string                `sorm:"description" json:"desc,omitempty"`
	Owner          string                `sorm:"owner" json:"owner,omitempty"`
	Achievements   string                `sorm:"achievements" json:"achievements,omitempty"`
	TypeOfOrg      int                   `sorm:"type_of_org" json:"type_of_org"`
	Followers      uint64                `sorm:"followers"  json:"follower_count"`
	Location       types.GeographyPoints `sorm:"location" json:"location"`
	Short          bool                  `sorm:"short,ignore" json:"-"`
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
