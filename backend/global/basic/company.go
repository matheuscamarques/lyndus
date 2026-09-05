package basic

import "bitbucket.org/lyndus/backend/infra/types"

type Company struct{
	ID                int        `json:"id" db:"id"`
	CNPJ              types.CNPJ `json:"cnpj,omitempty" db:"cnpj"`
	CompanyName       string     `json:"companyName" db:"company_name"`
	FantasyName       string     `json:"fantasyName" db:"fantasy_name"`
}