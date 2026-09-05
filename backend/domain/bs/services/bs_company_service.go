package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BSCompanyService struct {
	repo contracts.BSCompanyRepositoryInterface
}

var CompanyService contracts.BSCompanyServiceInterface

func NewBSCompanyService() *BSCompanyService {
	return &BSCompanyService{
		repo: repository.NewBSCompanyRepository(),
	}
}

func (bcs BSCompanyService) GetCompanyByBsID(bsID int) (company entity.Company, err error) {
	return bcs.repo.GetCompanyByBsID(bsID)
}
func (bcs BSCompanyService) UpdateCompanyByBS(company entity.Company) error {
	return bcs.repo.UpdateCompanyByBS(company)
}
func (bcs BSCompanyService) GetBSWeekDays(bsID int) (weekDays []entity.BMWeekDay, err error) {
	return bcs.repo.GetBSWeekDays(bsID)
}
func (bcs BSCompanyService) DeleteBSAllWeekDays(bsID int) error {
	return bcs.repo.DeleteBSAllWeekDays(bsID)
}
func (bcs BSCompanyService) AddBSWeekDay(bsID int, day entity.BMWeekDay) error {
	return bcs.repo.AddBSWeekDay(bsID, day)
}

func (bcs BSCompanyService) GetLocation(bsID int) (location string, err error) {
	return bcs.repo.GetLocation(bsID)
}
