package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BSCategoryService struct {
	repo contracts.BSCategoryRepositoryInterface
}

var CategoryService contracts.BSCategoryServiceInterface

func NewBSCategoryService() *BSCategoryService {
	return &BSCategoryService{
		repo: repository.NewBSCategoryRepository(),
	}
}

func (bcs BSCategoryService) GetBSCategories(bsID int) (categories []entity.Category, err error) {
	return bcs.repo.GetBSCategories(bsID)
}
func (bcs BSCategoryService) AddBSCategory(bsID, bsCategoryID int) error {
	return bcs.repo.AddBSCategory(bsID, bsCategoryID)
}
func (bcs BSCategoryService) DeleteBSCategory(bsID, bsCategoryID int) error {
	return bcs.repo.DeleteBSCategory(bsID, bsCategoryID)
}
