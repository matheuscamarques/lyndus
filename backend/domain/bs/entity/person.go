package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type Person struct {
	ID        int         `json:"id" db:"id"`
	Name      string      `json:"name" db:"name"`
	CPF       *types.CPF   `json:"cpf" db:"cpf"`
	Phone     string      `json:"phone" db:"phone"`
	Birthdate *types.Date `json:"birthdate" db:"birthdate"`

	BsID      int       `json:"BSID,omitempty" db:"bs_id"'`
	AppUserID *int      `json:"appUserID,omitempty" db:"app_user_id"`
	CreatedAT time.Time `json:"-" db:"created_at"`
}
