package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
	"bitbucket.org/lyndus/backend/infra/rest/response"
)

type ClientEmployeeService struct {
	repo contracts.ClientEmployeeRepositoryInterface
}

var EmployeeService contracts.ClientEmployeeServiceInterface

func NewClientEmployeeService() *ClientEmployeeService {
	return &ClientEmployeeService{
		repo: repository.NewClientEmployeeRepository(),
	}
}

func (ces ClientEmployeeService) GetEmployeeIDByCPF(employee *entity.Employee) error {
	return ces.repo.GetEmployeeIDByCPF(employee)
}

func (ces ClientEmployeeService) CreateEmployee(employee *entity.Employee) error {
	return ces.repo.CreateEmployee(employee)
}

func (ces ClientEmployeeService) UpdateEmployee(employee entity.Employee) error {
	return ces.repo.UpdateEmployee(employee)
}

func (ces ClientEmployeeService) UpdateEmployeeActive(employee entity.Employee) error {
	return ces.repo.UpdateEmployeeActive(employee)
}

func (ces ClientEmployeeService) GetEmployeeID(employee entity.Employee) (id response.ID, err error) {
	return ces.repo.GetEmployeeID(employee)
}

func (ces ClientEmployeeService) GetEmployee(employee *entity.Employee) error {
	return ces.repo.GetEmployee(employee)
}

func (ces ClientEmployeeService) DeleteEmployee(clientID, employeeID int) error {
	return ces.repo.DeleteEmployee(clientID, employeeID)
}

func (ces ClientEmployeeService) ActiveInactiveEmployee(clientID, employeeID int, active bool) error {
	return ces.repo.ActiveInactiveEmployee(clientID, employeeID, active)
}

func (ces ClientEmployeeService) GetEmployees(employee entity.Employee) (employees []entity.Employee, err error) {
	return ces.repo.GetEmployees(employee)
}

func (ces ClientEmployeeService) GetEmployeesCategoryValue(clientID int, active bool) (employees []entity.EmployeeBenefit, err error) {
	return ces.repo.GetEmployeesCategoryValue(clientID, active)
}

func (ces ClientEmployeeService) GetEmployeesByCategoryIDValue(clientID int, categories []int, active bool) (employees []entity.EmployeeBenefit, err error) {
	return ces.repo.GetEmployeesByCategoryIDValue(clientID, categories, active)
}

func (ces ClientEmployeeService) GetEmployeesByName(employee entity.Employee) (employees []entity.Employee, err error) {
	return ces.repo.GetEmployeesByName(employee)
}

func (ces ClientEmployeeService) AddEmployeeCategory(employee entity.Employee, companyCategoryID int) error {
	return ces.repo.AddEmployeeCategory(employee, companyCategoryID)
}

func (ces ClientEmployeeService) RemoveEmployeeCategory(employee entity.Employee, categoryID int) error {
	return ces.repo.RemoveEmployeeCategory(employee, categoryID)
}

func (ces ClientEmployeeService) GetEmployeeCategories(employeeID int) ([]entity.Category, error) {
	return ces.repo.GetEmployeeCategories(employeeID)
}

func (ces ClientEmployeeService) ChangeStatusEmployee(employee entity.Employee) error {
	return ces.repo.ChangeStatusEmployee(employee)
}

func (ces ClientEmployeeService) AddCategory(id int, category int) error {
	return ces.repo.AddCategory(id, category)
}
func (ces ClientEmployeeService) RemoveCategory(id int, category int) error {
	return ces.repo.RemoveCategory(id, category)
}

func (ces ClientEmployeeService) GetEmployeesV2(activePage, itemsPerPage, clientID int, search, orderBy string, sortDesc, active bool) (employees []entity.Employee, totalItems, totalPages int, err error) {
	if activePage == 0 {
		activePage = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 10
	} else if itemsPerPage > 100 {
		itemsPerPage = 100
	}

	employees, totalItems, totalPages, err = ces.repo.GetAll(
		activePage,
		itemsPerPage,
		clientID,
		search,
		orderBy,
		sortDesc,
		active)
	if employees == nil {
		employees = make([]entity.Employee, 0)
	}

	for k := range employees {
		employees[k].Categories, err = EmployeeService.GetEmployeeCategories(employees[k].ID)
		if err != nil {
			return
		}
		if employees[k].Categories == nil {
			employees[k].Categories = make([]entity.Category, 0)
		}
	}

	return employees, totalItems, totalPages, err
}
