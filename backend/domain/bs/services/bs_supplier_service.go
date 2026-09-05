package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/infra/criteria"
)

type BSSupplierService struct {
	repo contracts.BSSupplierRepositoryInterface
}

var SupplierService contracts.BSSupplierServiceInterface

func NewBSSupplierService() *BSSupplierService {
	return &BSSupplierService{
		repo: repository.NewBSSupplierRepository(),
	}
}

func (bss BSSupplierService) Create(supplier entity.Supplier) (id int, err error) {
	return bss.repo.Create(supplier)
}

func (bss BSSupplierService) Update(supplier entity.Supplier) error {
	return bss.repo.Update(supplier)
}

func (bss BSSupplierService) GetByID(bsID, supplierID int) (supplier entity.Supplier, err error) {
	return bss.repo.GetByID(bsID, supplierID)
}


func (bss BSSupplierService) GetAll(bsID, activePage, itemsPerPage int) (response criteria.BSSupplier, err error){


	criteria := criteria.Criteria{
		ActivePage:   activePage,
		ItemsPerPage: itemsPerPage,
	}

	err = criteria.Exec(bss.repo)

	if err != nil {
		return response, err
	}

	response.CResponse = response.New(criteria.ActivePage,criteria.TotalPages,criteria.TotalItems)
	response.Items, err = bss.repo.GetAll(bsID, criteria)

	if err != nil {
		return response, err
	}

	return response, err

}