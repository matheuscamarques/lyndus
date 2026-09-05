package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/db"
)

type ClientPollChoiceRepository struct {
	db.Connector
}

func NewClientPollChoiceRepository() (repo ClientPollChoiceRepository) {
	repo = ClientPollChoiceRepository{}
	repo.UsePostgres()
	repo.Table = "client_poll_choice"
	return repo
}

func (cpcr ClientPollChoiceRepository) Create(pollChoice entity.ClientPollChoice) (id int, err error) {
	command := `INSERT INTO 
							client_poll_choice(
												client_poll_question_id,
												choice_text
											 	)
				VALUES($1,$2) RETURNING id`
	stmt, err := cpcr.Conn.Prepare(command)
	if err != nil {
		return id, cpcr.Error(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		pollChoice.ClientPollQuestionID,
		pollChoice.ChoiceText,
	).Scan(&id)
	return id, cpcr.Error(err)
}
func (cpcr ClientPollChoiceRepository) Update(pollChoice entity.ClientPollChoice) (err error) {
	command := `UPDATE client_poll_choice 
								 SET choice_text=$3 
								 WHERE client_poll_question_id=$1
								 AND id=$2`
	stmt, err := cpcr.Conn.Prepare(command)
	if err != nil {
		return cpcr.Error(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(pollChoice.ClientPollQuestionID, pollChoice.ID, pollChoice.ChoiceText)
	return cpcr.Error(err)
}
func (cpcr ClientPollChoiceRepository) GetById(id int) (pollChoice entity.ClientPollChoice, err error) {
	command := `SELECT id, 
       					choice_text, 
       					client_poll_question_id 
       					 FROM client_poll_choice
						 WHERE  id = $1
						 `
	err = cpcr.Conn.Get(&pollChoice, command, id)
	return pollChoice, cpcr.Error(err)
}
