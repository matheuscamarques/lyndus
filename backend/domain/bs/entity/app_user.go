package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type AppUser struct {
	ID        int        `json:"id" db:"id"`
	PersonID  int        `json:"personID" db:"person_id"`
	Name      string     `json:"name" db:"name"`
	CPF       types.CPF  `json:"cpf" db:"cpf"`
	Cellphone string     `json:"phone" db:"cellphone"`
	Birthdate types.Date `json:"birthdate" db:"birthdate"`

	CreatedAT time.Time `json:"createdAT" db:"created_at"`
}
