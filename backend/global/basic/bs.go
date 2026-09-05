package basic

import "bitbucket.org/lyndus/backend/infra/types"

type BS struct {
	ID               int `json:"id" db:"id"`
	CNPJ types.CNPJ `json:"cnpj,omitempty" db:"cnpj"`
}
