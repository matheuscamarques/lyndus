package services

import (
	"bitbucket.org/lyndus/backend/domain/bonder/contracts"
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/domain/bonder/repository"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
	"log"
)

type BonderBsService struct {
	repo contracts.BonderBsRepositoryInterface
}

var BonderBS contracts.BonderBsServiceInterface

func NewBonderBsService() *BonderBsService {
	return &BonderBsService{
		repo: repository.NewBonderBsRepository(),
	}
}

// GetById(id int)
func (bbs BonderBsService) GetById(id int) (bs entity.BS, err error) {
	bs, err = bbs.repo.GetById(id)
	if err != nil {
		return bs, err
	}
	bs.BalanceDay, bs.RateAnticipation, bs.RateAnticipationBm, err = bbs.repo.GetLyndusPay(id)

	return bs, err
}

func (bbs BonderBsService) GetAll(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (bss []entity.BasicBS, totalItems, totalPages int, err error) {

	if activePage == 0 {
		activePage = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 10
	} else if itemsPerPage > 100 {
		itemsPerPage = 100
	}

	bss, totalItems, totalPages, err = bbs.repo.GetAll(
		activePage,
		itemsPerPage,
		search,
		orderBy,
		sortDesc,
		active)
	if bss == nil {
		bss = make([]entity.BasicBS, 0)
	}

	return bss, totalItems, totalPages, err
}

func (bbs BonderBsService) Register(bs entity.BS) (int, error) {
	return bbs.repo.Create(bs)
}

func (bbs BonderBsService) CreateBasicAccounts(bsID, closingDay int, rateAnticipation decimal.Decimal) error {
	return bbs.repo.CreateBasicAccounts(bsID, closingDay, rateAnticipation)
}

func (bbs BonderBsService) SearchByCNPJ(cnpj types.CNPJ) (id int, err error) {
	return bbs.repo.SearchByCNPJ(cnpj)
}

func (bbs BonderBsService) GetByIdSimple(id int) (bs entity.BS, err error) {
	return bbs.repo.GetByIdSimple(id)
}

func (bbs BonderBsService) Update(bs entity.BS) error {
	return bbs.repo.Update(bs)
}

func (bbs BonderBsService) UpdateLyndusActive(bsId int, active bool) (err error) {
	log.Println("s", bsId, active)
	return bbs.repo.UpdateLyndusActive(bsId, active)
}
func (bbs BonderBsService) Detete(bsId int) (err error) {
	return bbs.repo.Detete(bsId)
}

func (bbs BonderBsService) GetFinancialList(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (wcsList []entity.WithdrawalCashSimple, totalItems, totalPages int, err error) {
	return bbs.repo.GetFinancialList(activePage, itemsPerPage, search, orderBy, sortDesc, active)
}

func (bbs BonderBsService) GetBSWithdrawalCash(id int) (withdrawalCash entity.WithdrawalCash, err error) {
	return bbs.repo.GetBSWithdrawalCash(id)
}

func (bbs BonderBsService) GetBSWithdrawalCashBalance(bsID, id int) (wcbs []entity.WithdrawalCashBalance, err error) {
	return bbs.repo.GetBSWithdrawalCashBalance(bsID, id)
}

func (bbs BonderBsService) BSWithdrawalCashConfirm(bsID int, accountMovementID *int, accountReceiptID int, value, valueRate decimal.Decimal, wcbs []entity.WithdrawalCashBalance) (err error) {
	return bbs.repo.BSWithdrawalCashConfirm(bsID, accountMovementID, accountReceiptID, value, valueRate, wcbs)
}
