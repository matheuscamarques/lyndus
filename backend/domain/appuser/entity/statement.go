package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
)

type Statement struct {
	ID           int             `json:"id,omitempty" db:"id"`
	AppUserID    int             `json:"-" db:"app_user_id"`
	DateTime     types.DateTime  `json:"dateTime" db:"date_time"`
	Value        decimal.Decimal `json:"value" db:"value"`
	StatementsID int             `json:"-" db:"statements_id"`
	Name         string          `json:"name" db:"name"`
	Desc         string          `json:"desc,omitempty", db:"desc"`
}
