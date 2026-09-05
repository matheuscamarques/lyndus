package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID           int             `json:"id" db:"id"`
	BsID         int             `json:"bsID,omitempty" db:"bs_id"`
	BSPersonID   int             `json:"bsPersonID" db:"bs_person_id"`
	BSPersonName string          `json:"bsPersonName" db:"bs_person_name"`
	Value        decimal.Decimal `json:"value" db:"value"`
	Discount     decimal.Decimal `json:"discount" db:"discount"`
	TotalValue   decimal.Decimal `json:"totalValue" db:"total_value"`
	Obs          string          `json:"obs" db:"obs"`
	Items        []Item          `json:"items,omitempty" db:"items"`

	StatusID  int            `json:"statusID,omitempty" db:"order_status_id"`
	Status    string         `json:"status" db:"status"`
	CreatedAT types.DateTime `json:"createdAT" db:"created_at"`
}
