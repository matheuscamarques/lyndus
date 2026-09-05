package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
)

type AppUser struct {
	//TODO Validar
	AppUserPerson
	Balance decimal.NullDecimal `json:"balance" db:"balance"`

	FazBem    decimal.NullDecimal `json:"fazBem" db:"faz_bem"`
	LyndusBox decimal.NullDecimal `json:"lyndusBox" db:"lyndus_box"`
}

type Bonus struct {
	ID      int             `json:"id" db:"id"`
	Balance decimal.Decimal `json:"balance" db:"balance"`
}

type AppUserStatement struct {
	ID               int             `json:"id" db:"id"`
	AppUserID        int             `json:"appUserID" db:"app_user_id"`
	AppUserBalanceID int             `json:"appUserBalanceID" db:"app_user_balance_id"`
	StatementID      int             `json:"statementID" db:"statement_id"`
	DateTime         types.DateTime  `json:"dateTime" db:"date_time"`
	Value            decimal.Decimal `json:"value" db:"value"`
	Desc             string          `json:"desc" db:"desc"`
}

type AppUserPerson struct {
	ID                int        ` json:"id" db:"id"`
	Name              string     ` json:"name" db:"name"`
	Cpf               *types.CPF ` json:"cpf" db:"cpf"`
	Birthdate         types.Date ` json:"birthdate" db:"birthdate"`
	Phone             *string    ` json:"phone" db:"phone"`
	CellPhone         *string    ` json:"cellPhone" db:"cellphone"`
	Email             *string    ` json:"email,omitempty" db:"email"`
	State             *string    ` json:"state,omitempty" db:"state"`
	City              *string    ` json:"city,omitempty" db:"city"`
	District          *string    ` json:"district,omitempty" db:"district"`
	Street            *string    ` json:"street,omitempty" db:"street"`
	Number            *int       ` json:"number,omitempty" db:"number"`
	AddressComplement *string    ` json:"addressComplement,omitempty" db:"address_complement"`
	ZipCode           *string    ` json:"zipCode,omitempty" db:"zipcode"`
}
