package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/db"
)

type ClientPollCategoryRepository struct {
	db.Connector
}


func NewClientPollCategoryRepository() (repo ClientPollCategoryRepository) {
	repo.UsePostgres()
	repo.Table = "client_poll_category"
	return repo
}

func (cpcr ClientPollCategoryRepository) Create(pollCategory entity.ClientPollCategory) (id int, err error) {
	command := `INSERT INTO
								client_poll_category(
													 client_poll_id, client_category_id
												 	)
					VALUES($1,$2) RETURNING id`
	stmt, err := cpcr.Conn.Prepare(command)
	if err != nil {
		return id, cpcr.Error(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		pollCategory.ClientPollID,
		pollCategory.ClientCategoryID,
	).Scan(&id)
	return id, cpcr.Error(err)
}

func (cpcr ClientPollCategoryRepository) GetById(id int) (pollResults entity.ClientPollCategory, err error) {
	command := `SELECT id, 
      				   client_poll_id, 
       				   client_category_id FROM client_poll_category
						 WHERE  id = $1
						 `
	err = cpcr.Conn.Get(&pollResults, command, id)
	return pollResults, cpcr.Error(err)
}

func (cpcr ClientPollCategoryRepository) Update(pollCategory entity.ClientPollCategory) (err error) {
	return cpcr.Error(err)
}

func (cpcr ClientPollCategoryRepository) Exist(pollCategory entity.ClientPollCategory) (exist bool, err error) {
	command := `SELECT count(*)
					FROM client_poll_category 
					WHERE client_poll_id=$1 AND client_category_id=$2`
	stmt, err := cpcr.Conn.Prepare(command)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var total int
	err = stmt.QueryRow(pollCategory.ClientPollID, pollCategory.ClientCategoryID).Scan(&total)
	if err != nil {
		return false, err
	}
	return total > 0, err
}
