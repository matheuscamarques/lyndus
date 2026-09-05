package services

import (
	"bitbucket.org/lyndus/backend/domain/bonder/contracts"
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/domain/bonder/repository"
)

type BonderBSCategoryService struct {
	repo contracts.BonderBSCategoryRepositoryInterface
}

var BonderBSCategory contracts.BonderBSCategoryServiceInterface

func NewBonderBSCategoryService() *BonderBSCategoryService {
	return &BonderBSCategoryService{
		repo: repository.NewBonderBSCategoryRepository(),
	}
}

func (bbc BonderBSCategoryService) GetAll() (categories []entity.Category, err error) {
	return bbc.repo.GetAll()
}

func (bbc BonderBSCategoryService) Create(bsID, bsCategoryID int) error {
	return bbc.repo.Create(bsID, bsCategoryID)
}

func (bbc BonderBSCategoryService) GetByBS(bsID int) (categories []entity.Category, err error) {
	return bbc.repo.GetByBS(bsID)
}

func (bbc BonderBSCategoryService) Remove(bsID, bsCategoryID int) error {
	return bbc.repo.Remove(bsID, bsCategoryID)
}
