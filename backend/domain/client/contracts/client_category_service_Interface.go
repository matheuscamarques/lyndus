package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/rest/response"
)

type ClientCategoryServiceInterface interface {
	GetDefaultCategory(category *entity.Category) error
	CreateCategory(category entity.Category) (n int, err error)
	UpdateCategory(category entity.Category) error
	DeleteCategory(clientyID int, categoryID int) error
	ActiveInactiveCategory(clientID, categoryID int, active bool) error
	DeleteEmployeeCategory(companyCategoryID int) error
	GetCategoryIDByID(clientID int, categoryID int) (id response.ID, err error)
	GetCategoryByID(category *entity.Category) error
	GetCategoriesByName(clientID int, name string) (categories []entity.Category, err error)
	GetCategories(clientID int) (categories []entity.Category, err error)
	// GetCategoriesV2(clientID, page, itemsPerPage int, search string) (criteria.ClientCategoryResponse, error)
	GetCategoriesV2(activePage, itemsPerPage, clientID int, search, orderBy string, sortDesc, active bool) ([]entity.Category, int, int, error)
}
