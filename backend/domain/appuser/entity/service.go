package entity

import (
	"github.com/shopspring/decimal"
)

type Service struct {
	ID    int             `json:"id" db:"id"`
	Name  string          `json:"name" db:"name"`
	Desc  string          `json:"desc" db:"desc"`
	Value decimal.Decimal `json:"value,omitempty" db:"value"`
}

type ServiceCategory struct {
	ID       int       `json:"id" db:"id"`
	Desc     string    `json:"desc" db:"desc"`
	Services []Service `json:"services"`
}
