package entity

import "bitbucket.org/lyndus/backend/infra/types"

type Client struct {
	ID                int        `json:"id" db:"id"`
	CNPJ              types.CNPJ `json:"cnpj,omitempty" db:"cnpj"`
	CompanyName       string     `json:"companyName" db:"company_name"`
	AuthenticationID  int        `json:"-" db:"authentication_id"`
	CompanyID         int        `json:"companyID,omitempty" db:"company_id"`
	FantasyName       string     `json:"fantasyName" db:"fantasy_name"`
	Phone             string     `json:"phone,omitempty" db:"phone"`
	Email             string     `json:"email,omitempty" db:"email"`
	State             string     `json:"state,omitempty" db:"state"`
	City              string     `json:"city,omitempty" db:"city"`
	District          string     `json:"district,omitempty" db:"district"`
	Street            string     `json:"street,omitempty" db:"street"`
	Number            int        `json:"number,omitempty" db:"number"`
	AddressComplement string     `json:"addressComplement,omitempty" db:"address_complement"`
	Zipcode           string     `json:"zipcode,omitempty" db:"zipcode"`
	Desc              *string    `json:"desc,omitempty" db:"desc"`

	Active bool `db:"active_lyndus" json:"active"`
}
