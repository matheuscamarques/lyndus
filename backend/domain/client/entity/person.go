package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type Person struct {
	ID                int       `json:"id" db:"id"`
	GenderID          int       `json:"genderID" db:"gender_id"`
	Name              string    `json:"name" db:"name"`
	CPF               types.CPF `json:"cpf" db:"cpf"`
	Birthdate         time.Time `json:"birthdate" db:"birthdate"`
	Phone             string    `json:"phone" db:"phone"`
	Cellphone         string    `json:"cellphone" db:"cellphone"`
	Email             string    `json:"email" db:"email"`
	State             string    `json:"state" db:"state"`
	City              string    `json:"city" db:"city"`
	District          string    `json:"district" db:"district"`
	Street            string    `json:"street" db:"street"`
	Number            int       `json:"number" db:"number"`
	AddressComplement string    `json:"addressComplement" db:"address_complement"`
	Zipcode           string    `json:"zipcode" db:"zipcode"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at,omitempty"`
}

type PersonBasic struct {
	ID   int       `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
	CPF  types.CPF `json:"cpf" db:"cpf"`
}
