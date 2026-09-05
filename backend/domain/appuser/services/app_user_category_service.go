package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
)

type AppUserCategoryService struct {
	contracts.AppUserCategoryServiceInterface
	repo contracts.AppUserCategoryRepositoryInterface
}

var CategoryService contracts.AppUserCategoryServiceInterface

func NewAppUserCategoryService() *AppUserCategoryService {
	return &AppUserCategoryService{
		repo: repository.NewAppUserCategoryRepository(),
	}
}

func (aucs AppUserCategoryService) GetCategory(c *entity.Category) error {
	return aucs.repo.GetCategory(c)
}

func (aucs AppUserCategoryService) GetCategories() (categories []entity.Category, err error) {
	return aucs.repo.GetCategories()
}
