package repository

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/infra/db"
	"bitbucket.org/lyndus/backend/infra/types"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type BonderClientRepository struct {
	db.Connector
}

func NewBonderClientRepository() (repo BonderClientRepository) {
	repo.Table = "client"
	repo.UsePostgres()
	return repo
}

func (bcr BonderClientRepository) GetById(id int) (client entity.Client, err error) {
	command := `SELECT 
       					cli.id,
       					com.cnpj,
       					com.company_name,
       					com.fantasy_name,
       					com.phone,
       					com.email,
       					com.state,
       					com.city,
       					com.district,
       					com.street,
       					com.number,
       					com.address_complement,
       					com.zipcode,
       					com."desc",
       					cli.active_lyndus
					FROM client cli
					INNER JOIN company AS com ON cli.company_id = com.id
					WHERE cli.id = $1`
	err = bcr.Conn.Get(&client, command, id)
	return client, err
}

func (bcr BonderClientRepository) GetConnection() *sqlx.DB {
	return bcr.Conn
}

func (bcr BonderClientRepository) GetConnector() db.Connector {
	return bcr.Connector
}

func (bcr BonderClientRepository) GetTable() string {
	return bcr.Table
}

func (bcr BonderClientRepository) GetAll(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (companies []entity.Company, totalItems, totalPages int, err error) {
	command := `SELECT count(*)
					FROM client cli
					INNER JOIN company com ON com.id = cli.company_id
					WHERE cli.deleted = false
					  AND cli.active_lyndus = $1 
				`

	if search != "" {
		command = command + ` AND (unaccent(com.company_name) ILIKE $2 
							   OR unaccent(com.fantasy_name) ILIKE $2) `
		err = bcr.Conn.Get(&totalItems, command, active, "%"+search+"%")
	} else {
		err = bcr.Conn.Get(&totalItems, command, active)
	}

	if err != nil {
		return companies, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT 
					cli.id,
					com.company_name,
					com.fantasy_name,
					com.cnpj
					FROM client cli 
					INNER JOIN company  com ON cli.company_id = com.id
					WHERE cli.deleted = false
					  AND cli.active_lyndus = $3
					      `

	if search != "" {
		command = command + ` AND (unaccent(com.company_name) ILIKE $4 
							   OR unaccent(com.fantasy_name) ILIKE $4) `
	}

	if orderBy == "companyName" {
		command = command + ` ORDER BY com.company_name `
	} else if orderBy == "fantasyName" {
		command = command + ` ORDER BY com.fantasy_name `
	} else if orderBy == "cnpj" {
		command = command + ` ORDER BY com.cnpj `
	} else {
		command = command + ` ORDER BY com.company_name `
	}

	if sortDesc {
		command = command + " DESC "
	}
	command = command + ` OFFSET $1 LIMIT $2`

	if search != "" {
		err = bcr.Conn.Select(&companies, command, (activePage-1)*itemsPerPage, itemsPerPage, active, "%"+search+"%")
	} else {
		err = bcr.Conn.Select(&companies, command, (activePage-1)*itemsPerPage, itemsPerPage, active)
	}
	if err == sql.ErrNoRows {
		err = nil
	}

	return companies, totalItems, totalPages, err
}

func (bcr BonderClientRepository) Create(client entity.Client) (int, error) {
	// TODO
	command := `INSERT INTO client(authentication_id,company_id,cnpj) VALUES($1,$2,$3) RETURNING id`
	stmt, err := bcr.Conn.Prepare(command)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(client.AuthenticationID, client.CompanyID, client.CNPJ).Scan(&id)
	return id, err
}

func (bcr BonderClientRepository) SearchByCNPJ(cnpj types.CNPJ) (id int, err error) {
	command := `SELECT id FROM client WHERE cnpj = $1 `
	err = bcr.Conn.Get(&id, command, cnpj)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return 0, err
	}
	return id, err
}

func (bcr BonderClientRepository) GetCompanyId(id int) (idCompany int, err error) {
	command := `SELECT company_id FROM client WHERE id = $1`
	err = bcr.Conn.Get(&idCompany, command, id)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return idCompany, err
	}
	return idCompany, err
}

func (bcr BonderClientRepository) UpdateLyndusActive(clientId int, active bool) (err error) {
	command := `
		UPDATE client SET	
			active_lyndus=$2
		WHERE id = $1		
	`

	_, err = bcr.Conn.Exec(command, clientId, active)

	return err
}

func (bcr BonderClientRepository) Detete(clientID int) (err error) {
	command := `
		UPDATE client SET	
			deleted=$2
		WHERE id = $1		
	`
	_, err = bcr.Conn.Exec(command, clientID, true)

	return err
}
