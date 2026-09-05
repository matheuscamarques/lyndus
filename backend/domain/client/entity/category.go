package entity

import "github.com/shopspring/decimal"

type Category struct {
	ID       int             `json:"id,omitempty" db:"id"`
	ClientID int             `json:"clientID,omitempty" db:"client_id,omitempty"`
	Name     string          `json:"name,omitempty" db:"name"`
	Default  bool            `json:"default" db:"default"`
	Active   bool            `json:"active" db:"active"`
	Value    decimal.Decimal `json:"value" db:"value"`
}

type CategoryName struct {
	Category string `json:"category" db:"category"`
}
