package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/rest/response"
)

type ClientCategoryRepositoryInterface interface {
	GetDefaultCategory(category *entity.Category) error
	CreateCategory(category entity.Category) (n int, err error)
	UpdateCategory(category entity.Category) error
	DeleteCategory(clientID, categoryID int) error
	ActiveInactiveCategory(clientID, categoryID int, active bool) error
	DeleteEmployeeCategory(companyCategoryID int) error
	GetCategoryIDByID(clientID, categoryID int) (id response.ID, err error)
	GetCategoryByID(category *entity.Category) error
	GetCategoriesByName(clientID int, name string) (categories []entity.Category, err error)
	GetCategories(clientID int) (categories []entity.Category, err error)
	GetAll(activePage, itemsPerPage, clientID int, search, orderBy string, sortDesc, active bool) ([]entity.Category, int, int, error)
}
