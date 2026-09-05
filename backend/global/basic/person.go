package basic

import (
	"bitbucket.org/lyndus/backend/infra/types"
)

type Person struct {
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
