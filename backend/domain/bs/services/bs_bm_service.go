package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/domain/constants"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BSBMService struct {
	repo contracts.BSBMRepositoryInterface
}

var BmService contracts.BSBMServiceInterface

func NewBSBmService() *BSBMService {
	return &BSBMService{
		repo: repository.NewBSBMRepository(),
	}
}

func (bbs BSBMService) GetBMService(bmID, serviceID int) (entity.Service, error) {
	return bbs.repo.GetBMService(bmID, serviceID)
}

func (bbs BSBMService) GetBmID(bsID, bmID, bmStatusID int) (n int, err error) {
	return bbs.repo.GetBmID(bsID, bmID, bmStatusID)
}
func (bbs BSBMService) GetBM(bsID, bmID, bmStatusID int) (bm entity.BM, err error) {
	return bbs.repo.GetBM(bsID, bmID, bmStatusID)
}
func (bbs BSBMService) GetBmIDByCPF(bsID int, cpf types.CPF) (id int, err error) {
	return bbs.repo.GetBmIDByCPF(bsID, cpf)
}
func (bbs BSBMService) CreateBM(bm entity.BM) (id int, err error) {
	return bbs.repo.CreateBM(bm)
}
func (bbs BSBMService) UpdateBM(bm entity.BM) error {
	return bbs.repo.UpdateBM(bm)
}
func (bbs BSBMService) ChangeStatus(bsID, bmID, bmStatusID int) error {
	return bbs.repo.ChangeStatus(bsID, bmID, bmStatusID)
}
func (bbs BSBMService) AddBMService(bmID int, service entity.Service) error {
	return bbs.repo.AddBMService(bmID, service)
}
func (bbs BSBMService) DeleteBMService(bmID, serviceID int) error {
	return bbs.repo.DeleteBMService(bmID, serviceID)
}
func (bbs BSBMService) AddBMWeekDay(bmID int, day entity.BMWeekDay) error {
	return bbs.repo.AddBMWeekDay(bmID, day)
}
func (bbs BSBMService) DeleteBMAllWeekDays(bmID int) error {
	return bbs.repo.DeleteBMAllWeekDays(bmID)
}
func (bbs BSBMService) GetBMS(bsID, page, itemsPerPage int, search string) (bms []entity.BM, totalItems, totalPages int, err error) {

	if page == 0 {
		page = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 10
	} else if itemsPerPage > 100 {
		itemsPerPage = 100
	}

	if search != "" {
		var cpf types.CPF
		err = cpf.ParseCPF(search)

		if err == nil {
			bms, totalItems, totalPages, err = bbs.repo.GetBMSByCPF(page, itemsPerPage, bsID, constants.BMStatusActive, cpf)
			if err != nil {
				return
			}
		} else {
			bms, totalItems, totalPages, err = bbs.repo.GetBMSByName(page, itemsPerPage, bsID, constants.BMStatusActive, search)
			if err != nil {
				return
			}
		}
	} else {
		bms, totalItems, totalPages, err = bbs.repo.GetBMS(page, itemsPerPage, bsID, constants.BMStatusActive)
		if err != nil {
			return
		}
	}

	if bms == nil {
		bms = make([]entity.BM, 0)
	}
	return bms, totalItems, totalPages, err
}

func (bbs BSBMService) GetBmServicesIDs(bmID int) (services []int, err error) {
	return bbs.repo.GetBmServicesIDs(bmID)
}
func (bbs BSBMService) GetBmServices(bsID, bmID int) (services []entity.Service, err error) {
	return bbs.repo.GetBmServices(bsID, bmID)
}
func (bbs BSBMService) GetBMServices(bmID int) (services []entity.Service, err error) {
	return bbs.repo.GetBMServices(bmID)
}
func (bbs BSBMService) GetBMWeekDays(bmID int) (weekDays []entity.BMWeekDay, err error) {
	return bbs.repo.GetBMWeekDays(bmID)
}
func (bbs BSBMService) GetBMsByWeekDay(bsID, weekDay int) (weekDays []entity.BMBasic, err error) {
	return bbs.repo.GetBMsByWeekDay(bsID, weekDay)
}

func (bbs BSBMService) UpdateBMService(bmID int, service entity.Service) error {
	return bbs.repo.UpdateBMService(bmID, service)
}
