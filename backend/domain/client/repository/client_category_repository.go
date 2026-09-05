package repository

import (
	"bitbucket.org/lyndus/backend/infra/db"
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"github.com/jmoiron/sqlx"
)

type ClientCategoryRepository struct {
	db.Connector
}

func NewClientCategoryRepository() (repo ClientCategoryRepository) {
	repo.UsePostgres()
	repo.Table = "client_category"
	return repo
}

func (c ClientCategoryRepository) GetConnector() db.Connector {
	return c.Connector
}

func (c ClientCategoryRepository) GetTable() string {
	return c.Table
}

func (c ClientCategoryRepository) GetConnection() *sqlx.DB {
	return c.Conn
}

func (c ClientCategoryRepository) GetDefaultCategory(category *entity.Category) error {
	command := `SELECT cat.id,
					   cat.name
					FROM client_category cat
					WHERE cat.client_id = $1 AND cat.default = True`
	err := c.Conn.Get(category, command, category.ClientID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (c ClientCategoryRepository) CreateCategory(category entity.Category) (n int, err error) {
	command := `INSERT INTO client_category(client_id, name, value,"default") VALUES($1, $2, $3,$4) RETURNING ID`
	stmt, err := c.Conn.Prepare(command)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(category.ClientID, category.Name, category.Value, category.Default).Scan(&n)
	return
}

func (c ClientCategoryRepository) UpdateCategory(category entity.Category) error {
	command := `UPDATE client_category SET name=$3, value=$4 WHERE client_id=$1 AND id=$2`
	stmt, err := c.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(category.ClientID, category.ID, category.Name, category.Value)
	return err

}

func (c ClientCategoryRepository) DeleteCategory(clientID, categoryID int) error {
	command := `UPDATE client_category SET deleted=true WHERE client_id=$1 AND id=$2`
	stmt, err := c.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientID, categoryID)
	if err != nil {
		return err
	}
	command = `DELETE FROM client_employee_category WHERE client_category_id=$1`
	_, err = c.Conn.Exec(command, categoryID)
	return err
}

func (c ClientCategoryRepository) ActiveInactiveCategory(clientID, categoryID int, active bool) error {
	command := `UPDATE client_category SET active=$3 WHERE client_id=$1 AND id=$2`
	stmt, err := c.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientID, categoryID, active)
	return err
}

func (c ClientCategoryRepository) DeleteEmployeeCategory(companyCategoryID int) error {
	command := `DELETE FROM client_employee_category WHERE client_category_id=$1`
	stmt, err := c.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(companyCategoryID)
	return err

}

func (c ClientCategoryRepository) GetCategoryIDByID(clientID, categoryID int) (id response.ID, err error) {
	command := `SELECT cat.id
					FROM client_category cat
					WHERE cat.client_id = $1 AND cat.id= $2 AND (cat.deleted IS NULL OR cat.deleted = false) `
	err = c.Conn.Get(&id, command, clientID, categoryID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return id, err
}

func (c ClientCategoryRepository) GetCategoryByID(category *entity.Category) error {
	command := `SELECT cat.id,
					   cat.name,
					   cat.value,
       				   cat.active
					FROM client_category cat
					WHERE cat.client_id = $1
					  AND cat.id= $2
					  AND (cat.deleted IS NULL OR cat.deleted = false) `
	err := c.Conn.Get(category, command, category.ClientID, category.ID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (c ClientCategoryRepository) GetCategoriesByName(clientID int, name string) (categories []entity.Category, err error) {
	command := `SELECT id,
					   name,
					   value
					FROM client_category
					WHERE client_id =$1 AND lower(name) LIKE lower($2)
					AND (deleted IS NULL OR deleted = false)
					ORDER BY name `
	err = c.Conn.Select(&categories, command, clientID, name)
	if err == sql.ErrNoRows {
		err = nil
	}
	return categories, err
}

func (c ClientCategoryRepository) GetCategories(clientID int) (categories []entity.Category, err error) {
	command := `SELECT id,
					   name,
					   value
					FROM client_category
					WHERE client_id = $1
					  AND (deleted IS NULL OR deleted = false)
					  AND active = true
					ORDER BY name `
	err = c.Conn.Select(&categories, command, clientID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (c ClientCategoryRepository) GetAll(activePage, itemsPerPage, clientID int, search, orderBy string, sortDesc, active bool) (categories []entity.Category, totalItems, totalPages int, err error) {

	command := `SELECT count(*)
					FROM client_category
					WHERE (deleted IS NULL OR deleted = false)
					   AND active = $1
					   AND client_id = $2 `

	if search != "" {
		command = command + ` AND unaccent(name) ILIKE $3 `
		err = c.Conn.Get(&totalItems, command, active, clientID, "%"+search+"%")
	} else {
		err = c.Conn.Get(&totalItems, command, active, clientID)
	}

	if err != nil {
		return categories, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT id,
					  name,
					  value
					FROM client_category
					WHERE (deleted IS NULL OR deleted = false)
					   AND active = $3
					   AND client_id = $4 `

	if search != "" {
		command = command + ` AND unaccent(name) ILIKE $5 `
	}

	if orderBy == "name" {
		command = command + ` ORDER BY name `
	} else if orderBy == "value" {
		command = command + ` ORDER BY value `
	} else {
		command = command + ` ORDER BY name `
	}

	if sortDesc {
		command = command + " DESC "
	}
	command = command + ` OFFSET $1 LIMIT $2`

	if search != "" {
		err = c.Conn.Select(&categories, command, (activePage-1)*itemsPerPage, itemsPerPage, active, clientID, "%"+search+"%")
	} else {
		err = c.Conn.Select(&categories, command, (activePage-1)*itemsPerPage, itemsPerPage, active, clientID)
	}
	if err == sql.ErrNoRows {
		err = nil
	}

	return categories, totalItems, totalPages, err
}
