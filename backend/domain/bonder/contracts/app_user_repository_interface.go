package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
)

type AppUserRepositoryInterface interface {
	db.ConnectorInterface
	GetAppUserByCPF(cpf types.CPF, criteria criteria.Criteria) ([]aggregate.AppUserAggreate, error)
	GetAppUserByName(name string, criteria criteria.Criteria) ([]aggregate.AppUserAggreate, error)

	//GetAppUsers(criteria criteria.Criteria) (response criteria.AppUserResponse, err error)
	GetAppUsers(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (appUsers []entity.AppUserPerson, totalItems, totalPages int, err error)

	UpdateLyndusBox(bonus entity.Bonus) error

	GetAppUserByID(appUserID int) (appUsers entity.AppUser, err error)

	GetBalance(appUserID int) (bonus entity.Bonus, err error)
	InsertBalance(bonus entity.Bonus) (int, error)
	InsertBalanceHistory(appUserID int, oldValue, newValue decimal.Decimal) (int, error)
	UpdateBalance(appUserID int, value decimal.Decimal) error
	//InsertStatement(statement entity.AppUserStatement) (int, error)
}
