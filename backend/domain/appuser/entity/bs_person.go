package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
)

type BSPerson struct {
	ID        int        `json:"id" db:"id"`
	BsID      int        `json:"-" db:"bs+id"`
	AppUserID *int       `json:"appUserID" db:"app_user_id"`
	Name      string     `json:"name" db:"name"`
	CPF       types.CPF  `json:"cpf" db:"cpf"`
	Birthdate types.Date `json:"birthdate"db:"birthdate"`
	Phone     string     `json:"phone"db:"phone"`
}
