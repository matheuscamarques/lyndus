package entity

import "bitbucket.org/lyndus/backend/infra/types"

type Client struct {
	ID               int `json:"id" db:"id"`
	AuthenticationID int `json:"authenticationID,omitempty" db:"authentication_id"`
	CompanyID int `json:"companyID,omitempty" db:"company_id"`
	CNPJ types.CNPJ `json:"cnpj,omitempty" db:"cnpj"`
}

func (c Client)Verify() bool{
	return  c.AuthenticationID !=0 && c.CompanyID != 0 && len(c.CNPJ) >= 14
}