package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/db"
)

type ClientPollResultsRepository struct {
	db.Connector
}

func NewClientPollResultsRepository() (repo ClientPollResultsRepository) {
	repo = ClientPollResultsRepository{}
	repo.UsePostgres()
	repo.Table = "client_poll_results"
	return repo
}

func (cprr ClientPollResultsRepository) Create(pollResults entity.ClientPollResults) (id int, err error) {
	command := `INSERT INTO 
							client_poll_results(
												client_poll_question_id, 
							                    client_poll_choice_id, 
							                    client_employee_id, 
							                    responded
											 	)
				VALUES($1,$2,$3,$4) RETURNING id`
	stmt, err := cprr.Conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&id)
	return id, err
}
func (cprr ClientPollResultsRepository) Update(pollResults entity.ClientPollResults) (err error) {
	command := `UPDATE client_poll_results 
								 SET responded = $1
								 WHERE  id=$2`
	stmt, err := cprr.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		pollResults.Responded,
		pollResults.ID,
	)
	return err
}
func (cprr ClientPollResultsRepository) GetById(id int) (pollResults entity.ClientPollResults, err error) {
	command := `SELECT 
       					id, 
       					client_poll_question_id, 
       					client_poll_choice_id, 
       					client_employee_id, 
       					responded, 
       					created_at		
       FROM client_poll_results
						 WHERE  id = $1
						 `
	err = cprr.Conn.Get(&pollResults, command, id)
	return pollResults, err
}
