package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
)

type AppUserCreate struct {
	ID                int        `json:"-" db:"id"`
	PersonID          int        `json:"-" db:"person_id"`
	GenderID          int        `json:"genderID" db:"gender_id"`
	GenderName        string     `json:"genderName,omitempty" db:"desc"`
	Name              string     `json:"name" db:"name"`
	CPF               types.CPF  `json:"cpf" db:"cpf"`
	Birthdate         types.Date `json:"birthdate" db:"birthdate"`
	Phone             string     `json:"phone" db:"phone"`
	Cellphone         string     `json:"cellphone" db:"cellphone"`
	Email             string     `json:"email" db:"email"`
	State             string     `json:"state" db:"state"`
	City              string     `json:"city" db:"city"`
	District          string     `json:"district" db:"district"`
	Street            string     `json:"street" db:"street"`
	Number            int        `json:"number" db:"number"`
	AddressComplement string     `json:"addressComplement" db:"address_complement"`
	Zipcode           string     `json:"zipcode" db:"zipcode"`
	Password          string     `json:"password,omitempty" db:"password"`
	AcceptedTerm      bool       `json:"acceptedTerm,omitempty" db:"accepted_term"`
}
