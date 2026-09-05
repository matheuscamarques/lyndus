package entity

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	base "bitbucket.org/lyndus/backend/global/entity"
)

type BS struct {
	//todo ajeitar essa struct
	base.BS
	base.Company
	ID               int        `db:"id" json:"id"`
	CompanyName       string `json:"companyName" db:"company_name"`
	Phone             string `json:"phone" db:"phone"`
	FantasyName       string `json:"fantasyName" db:"fantasy_name"`
	Email             string `json:"email" db:"email"`
	State             string `json:"state" db:"state"`
	City              string `json:"city" db:"city"`
	District          string `json:"district" db:"district"`
	Street            string `json:"street" db:"street"`
	Number            int    `json:"number" db:"number"`
	AddressComplement string `json:"addressComplement" db:"address_complement"`
	Zipcode           string `json:"zipcode" db:"zipcode"`

	WeekDays        []aggregate.WeekDay `json:"weekDays,omitempty" db:"week_days"`
	Socials         []entity.BsSocial `json:"socials" db:"socials"`
}
