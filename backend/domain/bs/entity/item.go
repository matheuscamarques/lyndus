package entity

import (
	"github.com/shopspring/decimal"
)

type Item struct {
	ID          int             `json:"id" db:"id"`
	BSOrderID   int             `json:"-" db:"bs_order_id""`
	BSServiceID *int            `json:"serviceID" db:"bs_service_id"`
	BSProductID *int            `json:"productID" db:"bs_product_id"`
	BSBMID      int             `json:"bmID" db:"bs_bm_id"`
	BSBMName    string          `json:"bmName" db:"bs_bm_name"`
	Name        string          `json:"name" db:"name"`
	Value       decimal.Decimal `json:"value" db:"value"`
	Discount    decimal.Decimal `json:"discount" db:"discount"`
	TotalValue  decimal.Decimal `json:"totalValue" db:"total_value"`
	Quantity    int64           `json:"quantity" db:"quantity"`
}
