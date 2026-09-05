package entity

import "bitbucket.org/lyndus/backend/infra/types"

type Company struct {
	ID                int        `json:"id" db:"id"`
	CNPJ              types.CNPJ `json:"cnpj" db:"cnpj"`
	CompanyName       string     `json:"companyName" db:"company_name"`
	FantasyName       string     `json:"fantasyName" db:"fantasy_name"`
	Phone             string     `json:"phone" db:"phone"`
	Email             string     `json:"email" db:"email"`
	State             string     `json:"state" db:"state"`
	City              string     `json:"city" db:"city"`
	District          string     `json:"district" db:"district"`
	Street            string     `json:"street" db:"street"`
	Number            int        `json:"number" db:"number"`
	AddressComplement string     `json:"addressComplement" db:"address_complement"`
	Zipcode           string     `json:"zipcode" db:"zipcode"`
}
