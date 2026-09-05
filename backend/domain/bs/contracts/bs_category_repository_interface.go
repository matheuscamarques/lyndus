package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSCategoryServiceInterface interface {
	GetBSCategories(bsID int) (categories []entity.Category, err error)
	AddBSCategory(bsID, bsCategoryID int) error
	DeleteBSCategory(bsID, bsCategoryID int) error
}
