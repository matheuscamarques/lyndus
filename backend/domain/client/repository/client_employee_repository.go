package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/db"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ClientEmployeeRepository struct {
	db.Connector
}

func (e ClientEmployeeRepository) GetConnector() db.Connector {
	return e.Connector
}

func (e ClientEmployeeRepository) GetTable() string {
	return e.Table
}

func (e ClientEmployeeRepository) GetConnection() *sqlx.DB {
	return e.Conn
}

func NewClientEmployeeRepository() (repo ClientEmployeeRepository) {
	repo.UsePostgres()
	repo.Table = `client_employee`
	return repo
}

func (e ClientEmployeeRepository) GetEmployeeIDByCPF(employee *entity.Employee) error {
	command := `SELECT cem.id,
					   active 
					FROM client_employee cem
					WHERE cem.client_id = $1 AND cem.cpf = $2`
	err := e.Conn.Get(employee, command, employee.ClientID, employee.CPF)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (e ClientEmployeeRepository) CreateEmployee(employee *entity.Employee) error {
	command := `INSERT INTO client_employee(person_id, client_id, name, cpf)
					VALUES($1, $2, $3, $4) RETURNING ID`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(employee.PersonID, employee.ClientID, employee.Name, employee.CPF).Scan(&employee.ID)
	return err
}

func (e ClientEmployeeRepository) UpdateEmployee(employee entity.Employee) error {
	command := `UPDATE client_employee SET name=$3, cpf=$4 WHERE id=$2 AND client_id=$1`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	//log.Println(employee.ClientID, employee.ID, employee.Name, employee.CPF)
	_, err = stmt.Exec(employee.ClientID, employee.ID, employee.Name, employee.CPF)
	return err
}

func (e ClientEmployeeRepository) ActiveInactiveEmployee(clientID, employeeID int, active bool) error {
	command := `UPDATE client_employee SET active=$3 WHERE id=$2 AND client_id=$1`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientID, employeeID, active)
	return err
}

func (e ClientEmployeeRepository) UpdateEmployeeActive(employee entity.Employee) error {
	command := `UPDATE client_employee SET name=$3, active=$4 WHERE id=$2 AND client_id=$1`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(employee.ClientID, employee.ID, employee.Name, employee.Active)
	return err
}

func (e ClientEmployeeRepository) GetEmployeeID(employee entity.Employee) (id response.ID, err error) {
	command := `SELECT cem.id
					FROM client_employee cem
					WHERE cem.client_id = $1 AND cem.id= $2`
	err = e.Conn.Get(&id, command, employee.ClientID, employee.ID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return id, err
}

func (e ClientEmployeeRepository) GetEmployee(employee *entity.Employee) error {
	command := `SELECT cem.id,
					   cem.name,
					   cem.cpf,
       				   cem.active
					FROM client_employee cem
					WHERE cem.deleted = false
					  AND cem.client_id = $1
					  AND cem.id= $2 `
	err := e.Conn.Get(employee, command, employee.ClientID, employee.ID)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (e ClientEmployeeRepository) GetEmployees(employee entity.Employee) (employees []entity.Employee, err error) {
	command := `SELECT cem.id,
					   cem.name,
					   cem.cpf
					FROM client_employee cem
					WHERE cem.client_id= $1 AND cem.active = $2 ORDER BY cem.name `
	err = e.Conn.Select(&employees, command, employee.ClientID, employee.Active)
	if err == sql.ErrNoRows {
		err = nil
	}
	return employees, err
}

// func (a *Authentication) get_employees_from_category( companyID int, categoryID int):
//     return await db.fetch_rows("""SELECT cem.id, cem.name, cem.cpf FROM company_employee cem
//             INNER JOIN
//                 WHERE cem.company_id=$1""",
//                                company_id, category_id)

func (e ClientEmployeeRepository) GetEmployeesCategoryValue(clientID int, active bool) (employees []entity.EmployeeBenefit, err error) {
	command := `SELECT cem.id,
					   cem.name,
					   cem.cpf,
					   cat.id as category_id,
					   cat.name as category_name,
       				   cat.value as category_value,
       				   aus.id as app_user_id
					FROM client_employee cem
					INNER JOIN client_employee_category cec ON cem.id = cec.client_employee_id
					INNER JOIN client_category cat ON cat.id = cec.client_category_id
					INNER JOIN app_user aus ON aus.person_id = cem.person_id
					WHERE cem.client_id= $1 
					  AND cem.active = $2
					  AND cat.active = true
					  AND cat.deleted = false`
	err = e.Conn.Select(&employees, command, clientID, active)
	if err == sql.ErrNoRows {
		err = nil
	}
	return employees, err
}

func (e ClientEmployeeRepository) GetEmployeesByCategoryIDValue(clientID int, categories []int, active bool) (employees []entity.EmployeeBenefit, err error) {
	command := `SELECT cem.id,
					   cem.name,
					   cem.cpf,
					   cat.id as category_id,
					   cat.name as category_name,
       				   cat.value as category_value,
       				   aus.id as app_user_id
					FROM client_employee cem
					INNER JOIN client_employee_category cec ON cem.id = cec.client_employee_id
					INNER JOIN client_category cat on cat.id = cec.client_category_id
					INNER JOIN app_user aus ON aus.person_id = cem.person_id
					WHERE cem.client_id= $1
					  AND cem.active = $2
					  AND cat.id = ANY($3)
					  AND cat.active = true
					  AND cat.deleted = false`
	err = e.Conn.Select(&employees, command, clientID, active, pq.Array(categories))

	if err == sql.ErrNoRows {
		err = nil
	}
	return employees, err
}

func (e ClientEmployeeRepository) GetEmployeesByName(employee entity.Employee) (employees []entity.Employee, err error) {
	command := `SELECT cem.id,
					   cem.name,
					   cem.cpf
					FROM client_employee cem
					WHERE cem.client_id = $1 AND active = $2 AND lower(cem.name) LIKE lower($3)
					ORDER BY cem.name `
	err = e.Conn.Get(&employees, command, employee.ClientID, employee.Active, employee.Name)
	if err == sql.ErrNoRows {
		err = nil
	}
	return employees, err
}

func (e ClientEmployeeRepository) AddEmployeeCategory(employee entity.Employee, companyCategoryID int) error {
	command := `INSERT INTO client_employee_category(client_employee_id, client_category_id) VALUES($1, $2)`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(employee.ID, companyCategoryID)
	return err
}

//func (a *Authentication) UpdateEmployeeCategory(employeeID, companyCategoryID int) error {
//	command := `UPDATE company_employee_category SET company_category_id=$2 WHERE company_employee_id=$1`
//	stmt, err := e.conn.Prepare(command)
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//	_, err = stmt.Exec(employeeID, companyCategoryID)
//	return err
//}

func (e ClientEmployeeRepository) RemoveEmployeeCategory(employee entity.Employee, categoryID int) error {
	command := `DELETE FROM client_employee_category WHERE client_employee_id = $2 AND client_category_id = $3`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(employee.ID, categoryID)
	return err
}

func (e ClientEmployeeRepository) GetEmployeeCategories(employeeID int) (categories []entity.Category, err error) {
	command := `SELECT cca.id,
					  cca.name,
					  cca.value
					FROM client_category cca
					INNER JOIN client_employee_category cec ON cca.id = cec.client_category_id
					WHERE cec.client_employee_id = $1 
					  AND cca.deleted = false
					  AND cca.active = true`
	err = e.Conn.Select(&categories, command, employeeID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return categories, err
}

func (e ClientEmployeeRepository) ChangeStatusEmployee(employee entity.Employee) error {
	command := `UPDATE client_employee SET active = $3 WHERE id=$2 AND client_id=$1`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(employee.ClientID, employee.ID, employee.Active)
	return err
}

func (e ClientEmployeeRepository) DeleteEmployee(clientID, employeeID int) error {
	command := `UPDATE client_employee SET deleted = $3 WHERE id=$2 AND client_id=$1`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientID, employeeID, true)
	return err
}

func (e ClientEmployeeRepository) AddCategory(id int, category int) (err error) {
	command := `INSERT INTO 
							client_employee_category(
														client_employee_id, 
														client_category_id
							    ) 
							VALUES($1, $2)`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, category)
	return err
}

func (e ClientEmployeeRepository) RemoveCategory(id int, category int) error {
	command := `DELETE FROM client_employee_category 
                       WHERE client_employee_id=$1 
                       AND   client_category_id=$2`
	stmt, err := e.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, category)
	return err
}

func (e ClientEmployeeRepository) GetAll(activePage, itemsPerPage, clientID int, search, orderBy string, sortDesc, active bool) (employees []entity.Employee, totalItems, totalPages int, err error) {

	command := `SELECT count(*)
					FROM client_employee 
					WHERE deleted = false
					  AND active = $1 
					  AND client_id = $2 `
	if search != "" {
		command = command + ` AND unaccent(name) ILIKE $3 `
		err = e.Conn.Get(&totalItems, command, active, clientID, "%"+search+"%")
	} else {
		err = e.Conn.Get(&totalItems, command, active, clientID)
	}
	if err != nil {
		return employees, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT id,
					  name,
					  cpf
					FROM client_employee
					WHERE deleted = false
					  AND active = $3
					  AND client_id = $4 `

	if search != "" {
		command = command + ` AND unaccent(name) ILIKE $5 `
	}

	if orderBy == "name" {
		command = command + ` ORDER BY name `
	} else if orderBy == "cpf" {
		command = command + ` ORDER BY cpf `
	} else {
		command = command + ` ORDER BY name `
	}

	if sortDesc {
		command = command + " DESC "
	}
	command = command + ` OFFSET $1 LIMIT $2`

	if search != "" {
		err = e.Conn.Select(&employees, command, (activePage-1)*itemsPerPage, itemsPerPage, active, clientID, "%"+search+"%")
	} else {
		err = e.Conn.Select(&employees, command, (activePage-1)*itemsPerPage, itemsPerPage, active, clientID)
	}
	if err == sql.ErrNoRows {
		err = nil
	}

	return employees, totalItems, totalPages, err
}

//func (e ClientEmployeeRepository) GetAll(ctra criteria.Criteria) (response criteria.ClientEmployeeResponse,err error) {
//	command := `SELECT client_employee.id,
//					   client_employee.name,
//					   client_employee.cpf
//					FROM client_employee
//					WHERE client_employee.client_id = $1 AND employee_status_id = $2 AND lower(client_employee::text) LIKE '%' || lower($3) || '%'
//					ORDER BY client_employee.name `
//	command = ctra.Query(command)
//	err = e.Conn.Select(&response.Items,command,ctra.Args...)
//	if err != nil {
//		return response, err
//	}
//	err = ctra.ExecWithQuery(e,`WHERE client_employee.client_id = $1 AND employee_status_id = $2 AND lower(client_employee::text) LIKE '%' || lower($3) || '%'`,ctra.Args...)
//	if err != nil {
//		return response, err
//	}
//	response.CResponse = response.NewWithCriteria(ctra)
//	if response.Items == nil {
//		response.Items = make([]entity.Employee, 0)
//	}
//
//	return response, err
//}
