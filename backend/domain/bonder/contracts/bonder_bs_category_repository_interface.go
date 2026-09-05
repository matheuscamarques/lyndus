package contracts

import "bitbucket.org/lyndus/backend/domain/bonder/entity"

type BonderBSCategoryRepositoryInterface interface {
	GetAll() (categories []entity.Category, err error)
	Create(bsID, bsCategoryID int) error
	GetByBS(bsID int) (categories []entity.Category, err error)
	Remove(bsID, bsCategoryID int) error
}
