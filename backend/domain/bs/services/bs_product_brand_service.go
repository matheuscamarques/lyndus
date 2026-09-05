package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/infra/criteria"
)

type BSProductBrandService struct {
	repo contracts.BSProductBrandRepositoryInterface
}


var ProductBrandService contracts.BSProductBrandServiceInterface

func NewBSProductBrandService() *BSProductBrandService {
	return &BSProductBrandService{
		repo: repository.NewBSProductBrandRepository(),
	}
}

func (bps BSProductBrandService) Create(brand entity.ProductBrand) (int, error) {
	return bps.repo.Create(brand)
}
func (bps BSProductBrandService) GetByID(brandID int) (entity.ProductBrand, error) {
	return bps.repo.GetByID(brandID)
}
func (bps BSProductBrandService) GetByName(name string)(brands entity.ProductBrand, err error){
	return bps.repo.GetByName(name)
}

func (bps BSProductBrandService) SearchByName(name string) (brands []entity.ProductBrand, err error) {
	return bps.repo.SearchByName(name)
}

func (bps BSProductBrandService) GetAll(activePage, itemsPerPag int) (response criteria.BSProductBrand, err error) {

	criteria := criteria.Criteria{
		ActivePage:   activePage,
		ItemsPerPage: itemsPerPag,
	}

	err = criteria.Exec(bps.repo)
	if err != nil {
		return response, err
	}

	response.CResponse = response.New(criteria.ActivePage,criteria.TotalPages,criteria.TotalItems)
	response.Items, err = bps.repo.GetAll(criteria)

	if err != nil {
		return response, err
	}

	return response, err
}
