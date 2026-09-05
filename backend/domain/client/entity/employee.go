package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
)

type Employee struct {
	ID         int        `json:"id"         db:"id"`
	Name       string     `json:"name"       db:"name"`
	CPF        types.CPF  `json:"cpf"        db:"cpf"`
	Categories []Category `json:"categories" db:"categories"`

	ClientID int `json:"clientID,omitempty" db:"-"`
	PersonID int `json:"personID,omitempty" db:"person_id,omitempty"`
	//StatusID int  `json:"statusID,omitempty" db:"status_id,omitempty"`
	Active bool `json:"active" db:"active"`
}

type EmployeeCreate struct {
	Name       string    `json:"name" db:"name"`
	CPF        types.CPF `json:"cpf" db:"cpf"`
	Categories []int     `json:"categories" db:"categories"`
}

type EmployeeBenefit struct {
	ID            int             `json:"id" db:"id"`
	Name          string          `json:"name" db:"name"`
	CPF           types.CPF       `json:"cpf" db:"cpf"`
	CategoryID    int             `json:"categoryID" db:"category_id"`
	CategoryName  string          `json:"categoryName" db:"category_name"`
	CategoryValue decimal.Decimal `json:"categoryValue" db:"category_value"`
	AppUserID     int             `json:"appUserID" db:"app_user_id"`
}
