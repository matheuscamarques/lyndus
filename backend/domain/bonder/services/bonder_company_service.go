package services

import (
	"bitbucket.org/lyndus/backend/domain/bonder/contracts"
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/domain/bonder/repository"
)

var BonderCompany contracts.BonderCompanyServiceInterface

type BonderCompanyService struct {
	repo contracts.BonderCompanyRepositoryInterface
}

func NewBonderCompanyService() *BonderCompanyService {
	return &BonderCompanyService{
		repo: repository.NewBonderCompanyRepository(),
	}
}

func (bcs BonderCompanyService) Register(company *entity.BS) (companyID int, err error) {
	//companyID, err = bcs.repo.GetIdByCNPJ(company.CNPJ)
	if err != nil {
		return 0, err
	}

	return bcs.repo.Create(company)
}

func (bcs BonderCompanyService) Update(company *entity.BS) error {
	return bcs.repo.Update(company)
}
