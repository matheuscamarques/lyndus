package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type Supplier struct {
	ID             int        `json:"id" db:"id"`
	BsID           int        `json:"-" db:"bs_id"`
	Name           string     `json:"name" db:"name"`
	CNPJ           types.CNPJ `json:"cnpj" db:"cnpj"`
	Phone          string     `json:"phone,omitempty" db:"phone"`
	Representative string     `json:"representative,omitempty" db:"representative"`
	Obs            string     `json:"obs,omitempty" db:"obs"`
	Address        string     `json:"address,omitempty" db:"address"`
	CreatedBy      int        `json:"createdBy,omitempty" db:"created_by"`
	CreatedAt      *time.Time  `json:"createdAt,omitempty" db:"created_at"`
}
