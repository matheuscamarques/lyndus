package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
	"github.com/shopspring/decimal"
)

type AppUserBalanceService struct {
	repo contracts.AppUserBalanceRepositoryInterface
}

var BalanceService contracts.AppUserBalanceServiceInterface

func NewAppUserBalanceService() *AppUserBalanceService {
	return &AppUserBalanceService{
		repo: repository.NewAppUserBalanceRepository(),
	}
}

func (aub AppUserBalanceService) GetBalance(appUserID int) (balance entity.Balance, err error) {
	return aub.repo.GetBalance(appUserID)
}

func (aub AppUserBalanceService) UpdateBalance(appUserID int, value decimal.Decimal) error {
	return aub.repo.UpdateBalance(appUserID, value)
}

func (aub AppUserBalanceService) InsertBalanceHistory(appUserID int, oldValue, newValue decimal.Decimal) (int, error) {
	return aub.repo.InsertBalanceHistory(appUserID, oldValue, newValue)
}
