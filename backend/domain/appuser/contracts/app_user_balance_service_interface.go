package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"github.com/shopspring/decimal"
)

type AppUserBalanceServiceInterface interface {
	GetBalance(appUserID int) (balance entity.Balance, err error)

	UpdateBalance(appUserID int, value decimal.Decimal) error
	InsertBalanceHistory(appUserID int, oldValue, newValue decimal.Decimal) (int, error)
}
