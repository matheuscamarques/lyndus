package contracts

import "bitbucket.org/lyndus/backend/domain/appuser/entity"

type AppUserCategoryServiceInterface interface {
	GetCategory(c *entity.Category) error
	GetCategories() (categories []entity.Category, err error)
}
