package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/global/contracts"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
)

type BonderBsRepositoryInterface interface {
	contracts.RepositoryInterface
	Create(bs entity.BS) (int, error)
	CreateBasicAccounts(bsID, closingDay int, rateAnticipation decimal.Decimal) error
	SearchByCNPJ(cnpj types.CNPJ) (id int, err error)
	//GetAll(criteria criteria.Criteria) (response criteria.BSResponse, err error)
	GetAll(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (bss []entity.BasicBS, totalItems, totalPages int, err error)
	GetById(id int) (entity.BS, error)
	GetLyndusPay(bsID int) (closingDay int, rateAnticipation, rateAnticipationBm decimal.Decimal, err error)

	UpdateLyndusActive(bsId int, active bool) (err error)
	Detete(bsId int) (err error)

	GetByIdSimple(id int) (bs entity.BS, err error)
	Update(bs entity.BS) error

	GetFinancialList(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) ([]entity.WithdrawalCashSimple, int, int, error)
	GetBSWithdrawalCash(id int) (withdrawalCash entity.WithdrawalCash, err error)
	GetBSWithdrawalCashBalance(bsID, id int) (wcbs []entity.WithdrawalCashBalance, err error)
	BSWithdrawalCashConfirm(bsID int, accountMovementID *int, accountReceiptID int, value, valueRate decimal.Decimal, wcbs []entity.WithdrawalCashBalance) (err error)
	//GetAllWithWhere(c criteria.Criteria, s string) ([]aggregate.BSAggregate, error)
}
