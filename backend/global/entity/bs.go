package entity

import "bitbucket.org/lyndus/backend/infra/types"

type BS struct {
	ID               int        `db:"id" json:"id,omitempty"`
	AuthenticationID int        `db:"authentication_id" json:"authenticationID,omitempty"`
	CompanyID        int        `db:"company_id" json:"companyID,omitempty"`
	Balance          int        `db:"balance" json:"balance,omitempty"`
	Lat              *float64   `db:"lat" json:"lat,omitempty"`
	Lon              *float64   `db:"lon" json:"lon,omitempty"`
	Distance         *float64   `db:"distance" json:"distance,omitempty"`
	CNPJ             types.CNPJ `db:"cnpj" json:"cnpj,omitempty"`
}
