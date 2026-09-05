package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
	"bitbucket.org/lyndus/backend/infra/rest/response"
)

type ClientCategoryService struct {
	Repo contracts.ClientCategoryRepositoryInterface
}

var CategoryService contracts.ClientCategoryServiceInterface

func NewClientCategoryService() *ClientCategoryService {
	return &ClientCategoryService{
		Repo: repository.NewClientCategoryRepository(),
	}
}

func (ccs ClientCategoryService) GetDefaultCategory(category *entity.Category) error {
	return ccs.Repo.GetDefaultCategory(category)
}
func (ccs ClientCategoryService) CreateCategory(category entity.Category) (n int, err error) {
	return ccs.Repo.CreateCategory(category)
}
func (ccs ClientCategoryService) UpdateCategory(category entity.Category) error {
	return ccs.Repo.UpdateCategory(category)
}
func (ccs ClientCategoryService) DeleteCategory(companyID int, categoryID int) error {
	return ccs.Repo.DeleteCategory(companyID, categoryID)
}
func (ccs ClientCategoryService) ActiveInactiveCategory(clientID, categoryID int, active bool) error {
	return ccs.Repo.ActiveInactiveCategory(clientID, categoryID, active)
}

func (ccs ClientCategoryService) DeleteEmployeeCategory(companyCategoryID int) error {
	return ccs.Repo.DeleteEmployeeCategory(companyCategoryID)
}
func (ccs ClientCategoryService) GetCategoryIDByID(clientID int, categoryID int) (id response.ID, err error) {
	return ccs.Repo.GetCategoryIDByID(clientID, categoryID)
}
func (ccs ClientCategoryService) GetCategoryByID(category *entity.Category) error {
	return ccs.Repo.GetCategoryByID(category)
}
func (ccs ClientCategoryService) GetCategoriesByName(clientID int, name string) (categories []entity.Category, err error) {
	return ccs.Repo.GetCategoriesByName(clientID, name)
}
func (ccs ClientCategoryService) GetCategories(clientID int) (categories []entity.Category, err error) {
	return ccs.Repo.GetCategories(clientID)
}

func (ccs ClientCategoryService) GetCategoriesV2(activePage, itemsPerPage, clientID int, search, orderBy string, sortDesc, active bool) (categories []entity.Category, totalItems, totalPages int, err error) {

	if activePage == 0 {
		activePage = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 10
	} else if itemsPerPage > 100 {
		itemsPerPage = 100
	}

	categories, totalItems, totalPages, err = ccs.Repo.GetAll(
		activePage,
		itemsPerPage,
		clientID,
		search,
		orderBy,
		sortDesc,
		active)
	if categories == nil {
		categories = make([]entity.Category, 0)
	}

	return categories, totalItems, totalPages, err
}
