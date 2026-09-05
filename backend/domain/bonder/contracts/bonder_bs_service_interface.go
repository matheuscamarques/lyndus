package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
)

type BonderBsServiceInterface interface {
	Register(bs entity.BS) (int, error)
	//Edit(bs entity.BS) (int,error)
	SearchByCNPJ(cnpj types.CNPJ) (int, error)

	CreateBasicAccounts(bsID, closingDay int, rateAnticipation decimal.Decimal) error
	// GetAll(ActivePage, ItemsPerPage int,search string) (criteria.BSResponse, error)

	GetAll(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (bss []entity.BasicBS, totalItems, totalPages int, err error)
	GetById(id int) (entity.BS, error)

	UpdateLyndusActive(bsId int, active bool) (err error)
	Detete(bsId int) (err error)

	GetByIdSimple(id int) (bs entity.BS, err error)
	Update(bs entity.BS) error

	GetFinancialList(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (wcsList []entity.WithdrawalCashSimple, totalItems, totalPages int, err error)
	GetBSWithdrawalCash(id int) (withdrawalCash entity.WithdrawalCash, err error)
	GetBSWithdrawalCashBalance(bsID, id int) (wcbs []entity.WithdrawalCashBalance, err error)
	BSWithdrawalCashConfirm(bsID int, accountMovementID *int, accountReceiptID int, value, valueRate decimal.Decimal, wcbs []entity.WithdrawalCashBalance) (err error)
}
