package entity

import (
	"github.com/shopspring/decimal"
)

type Service struct {
	ID           int             `json:"id" db:"id"`
	BsID         int             `json:"bsID,omitempty" db:"bs_id"`
	BsBmID       int             `json:"bsBmID,omitempty" db:"bs_bm_id"`
	CategoryID   int             `json:"categoryID,omitempty" db:"bs_service_category_id"`
	CategoryName string          `json:"categoryName,omitempty" db:"category_name"`
	Name         string          `json:"name" bd:"name"`
	Desc         string          `json:"desc" db:"desc"`
	Value        decimal.Decimal `json:"value" db:"value"`
	Duration     int             `json:"duration,omitempty" db:"duration"`
}

type ServiceCategory struct {
	ID   int    `json:"id" db:"id"`
	Desc string `json:"desc" db:"desc"`
}
