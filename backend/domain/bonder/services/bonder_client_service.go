package services

import (
	"bitbucket.org/lyndus/backend/domain/bonder/contracts"
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/domain/bonder/repository"
	"bitbucket.org/lyndus/backend/infra/types"
)

var BonderClient contracts.BonderClientServiceInterface

type BonderClientService struct {
	repo contracts.BonderClientRepositoryInterface
}

func NewBonderClientService() *BonderClientService {
	return &BonderClientService{
		repo: repository.NewBonderClientRepository(),
	}
}

func (bcs BonderClientService) GetById(id int) (client entity.Client, err error) {
	return bcs.repo.GetById(id)
}

func (bcs BonderClientService) GetAll(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (companies []entity.Company, totalItems, totalPages int, err error) {
	if activePage == 0 {
		activePage = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 10
	} else if itemsPerPage > 100 {
		itemsPerPage = 100
	}

	companies, totalItems, totalPages, err = bcs.repo.GetAll(activePage, itemsPerPage, search, orderBy, sortDesc, active)
	if companies == nil {
		companies = make([]entity.Company, 0)
	}

	return companies, totalItems, totalPages, err
}

func (bcs BonderClientService) Register(client entity.Client) (int, error) {
	return bcs.repo.Create(client)
}

func (bcs BonderClientService) SearchByCNPJ(cnpj types.CNPJ) (id int, err error) {
	return bcs.repo.SearchByCNPJ(cnpj)
}

func (bcs BonderClientService) GetCompanyId(id int) (int, error) {
	return bcs.repo.GetCompanyId(id)
}
func (bcs BonderClientService) UpdateLyndusActive(clientId int, active bool) (err error) {
	return bcs.repo.UpdateLyndusActive(clientId, active)
}
func (bcs BonderClientService) Detete(clientId int) (err error) {
	return bcs.repo.Detete(clientId)
}
