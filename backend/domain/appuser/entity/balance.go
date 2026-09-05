package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Balance struct {
	ID         int             `json:"-" db:"id"`
	AppUserID  int             `json:"-" db:"app_user_id"`
	Value      decimal.Decimal `json:"value" db:"value"`
	LastUpdate time.Time       `json:"-" db:"last_update"`
}
