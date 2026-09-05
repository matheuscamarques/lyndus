package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
)

type ClientCompanyService struct {
	contracts.ClientCompanyServiceInterface
	repo contracts.ClientCompanyRepositoryInterface
}

var CompanyService contracts.ClientCompanyServiceInterface

func NewClientCompanyService() *ClientCompanyService {
	return &ClientCompanyService{
		repo: repository.NewClientCompanyRepository(),
	}
}

func (ccs ClientCompanyService) GetCompanyByID(targetID int) (entity.Company, error) {
	return ccs.repo.GetCompanyByID(targetID)
}

func (ccs ClientCompanyService) UpdateCompany(company entity.Company) error {
	return ccs.repo.UpdateCompany(company)
}
