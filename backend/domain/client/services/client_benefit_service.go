package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
)

type ClientBenefitService struct {
	repo contracts.ClientBenefitRepositoryInterface
}

var BenefitService contracts.ClientBenefitServiceInterface

func NewClientBenefitService() *ClientBenefitService {
	return &ClientBenefitService{
		repo: repository.NewClientBenefitRepository(),
	}
}

func (cbs ClientBenefitService) CreateBenefitUserCategory(buc *entity.BenefitUserCategory) (int, error) {
	return cbs.repo.CreateBenefitUserCategory(buc)
}

func (cbs ClientBenefitService) CreateBenefit(b *entity.Benefit) error {
	return cbs.repo.CreateBenefit(b)
}
func (cbs ClientBenefitService) CreateBenefitUser(bu entity.BenefitUser) (int, error) {
	return cbs.repo.CreateBenefitUser(bu)
}
func (cbs ClientBenefitService) UpdateBenefitValue(b *entity.Benefit) error {
	return cbs.repo.UpdateBenefitValue(b)
}
func (cbs ClientBenefitService) UpdateBenefitUserValue(bu *entity.BenefitUser) error {
	return cbs.repo.UpdateBenefitUserValue(bu)
}

func (cbs ClientBenefitService) UpdateBenefitUserAdditional(bu *entity.BenefitUser) error {
	return cbs.repo.UpdateBenefitUserAdditional(bu)
}
func (cbs ClientBenefitService) GetBenefitStatus(b *entity.Benefit) (status int, err error) {
	return cbs.repo.GetBenefitStatus(b)
}
func (cbs ClientBenefitService) GetBenefit(b entity.Benefit) (benefitResp entity.BenefitRespFull, err error) {
	return cbs.repo.GetBenefit(b)
}
func (cbs ClientBenefitService) UpdateBenefitDesc(b entity.Benefit) error {
	return cbs.repo.UpdateBenefitDesc(b)
}
func (cbs ClientBenefitService) UpdateBenefitStatus(b entity.Benefit) error {
	return cbs.repo.UpdateBenefitStatus(b)
}
func (cbs ClientBenefitService) GetBenefitByStatus(b entity.Benefit) error {
	return cbs.repo.GetBenefitByStatus(b)
}
func (cbs ClientBenefitService) GetBenefits(clientID int) (benefitsResp []entity.BenefitResp, err error) {
	return cbs.repo.GetBenefits(clientID)
}
func (cbs ClientBenefitService) GetBenefitUser(b entity.Benefit) (benefitUser entity.BenefitUserResp, err error) {
	return cbs.repo.GetBenefitUser(b)
}
func (cbs ClientBenefitService) GetBenefitsUsers(b entity.Benefit) (benefitsUser []entity.BenefitUserResp, err error) {
	return cbs.repo.GetBenefitsUsers(b)
}
func (cbs ClientBenefitService) GetBenefitsUserCategory(bu entity.BenefitUserResp) (categories []string, err error) {
	return cbs.repo.GetBenefitsUserCategory(bu)
}
func (cbs ClientBenefitService) GetBenefitsUsersValues(bu entity.BenefitUser) (benefitsUser []entity.BenefitUser, err error) {
	return cbs.repo.GetBenefitsUsersValues(bu)
}

func (cbs ClientBenefitService) GetBenefitsV2(activePage, itemsPerPage, clientID, statusID int, search, orderBy string, sortDesc bool) (benefits []entity.BenefitResp, totalItems, totalPages int, err error) {
	if activePage == 0 {
		activePage = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 10
	} else if itemsPerPage > 100 {
		itemsPerPage = 100
	}

	benefits, totalItems, totalPages, err = cbs.repo.GetAll(
		activePage,
		itemsPerPage,
		clientID,
		statusID,
		search,
		orderBy,
		sortDesc)
	if benefits == nil {
		benefits = make([]entity.BenefitResp, 0)
	}
	return benefits, totalItems, totalPages, err
}
