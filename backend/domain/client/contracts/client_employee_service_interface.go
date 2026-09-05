package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/rest/response"
)

type ClientEmployeeServiceInterface interface {
	GetEmployeeIDByCPF(employee *entity.Employee) error
	CreateEmployee(employee *entity.Employee) error
	UpdateEmployee(employee entity.Employee) error
	UpdateEmployeeActive(employee entity.Employee) error
	ActiveInactiveEmployee(clientID, employeeID int, active bool) error
	GetEmployeeID(employee entity.Employee) (id response.ID, err error)
	DeleteEmployee(clientID, employeeID int) error
	GetEmployee(employee *entity.Employee) error
	GetEmployees(employee entity.Employee) (employees []entity.Employee, err error)
	GetEmployeesCategoryValue(clientID int, active bool) (employees []entity.EmployeeBenefit, err error)
	GetEmployeesByCategoryIDValue(clientID int, categories []int, active bool) (employees []entity.EmployeeBenefit, err error)
	GetEmployeesByName(employee entity.Employee) (employees []entity.Employee, err error)
	AddEmployeeCategory(employee entity.Employee, companyCategoryID int) error
	RemoveEmployeeCategory(employee entity.Employee, categoryID int) error
	GetEmployeeCategories(employeeID int) (categories []entity.Category, err error)
	ChangeStatusEmployee(employee entity.Employee) error
	AddCategory(id int, category int) error
	RemoveCategory(id int, category int) error
	GetEmployeesV2(activePage, itemsPerPage, clientID int, search, orderBy string, sortDesc, active bool) (employees []entity.Employee, totalItems, totalPages int, err error)
}
