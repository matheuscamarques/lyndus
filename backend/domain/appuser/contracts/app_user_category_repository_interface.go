package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
)

type AppUserCategoryRepositoryInterface interface {
	GetCategory(c *entity.Category) error
	GetCategories() (categories []entity.Category, err error)
}
