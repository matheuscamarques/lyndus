package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/global/basic"
	"bitbucket.org/lyndus/backend/infra/db"
	"fmt"
)

type ClientPollQuestionRepository struct {
	db.Connector
}



func NewClientPollQuestionRepository() (repo ClientPollQuestionRepository) {
	repo.UsePostgres()
	repo.Table = "client_poll_question"
	return repo
}

func (cpqr ClientPollQuestionRepository) Create(pollQuestion entity.ClientPollQuestion) (id int, err error) {
	command := `INSERT INTO 
							client_poll_question(
										client_poll_id, 
							            multiple_choice, 
							            question_text
											 	)
				VALUES($1,$2,$3) RETURNING id`
	stmt, err := cpqr.Conn.Prepare(command)
	if err != nil {
		return id, cpqr.Error(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		pollQuestion.ClientPollID,
		pollQuestion.MultipleChoice,
		pollQuestion.QuestionText,
	).Scan(&id)
	return id, cpqr.Error(err)
}
func (cpqr ClientPollQuestionRepository) Update(pollQuestion entity.ClientPollQuestion) (err error) {
	command := `UPDATE client_poll_question 
								 SET 
								     question_text = $2,
								     multiple_choice = $3
								 WHERE id=$1`
	fmt.Println("AQUI 2",pollQuestion.ID)
	stmt, err := cpqr.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		pollQuestion.ID,
		pollQuestion.QuestionText,
		pollQuestion.MultipleChoice,
	)
	return err
}
func (cpqr ClientPollQuestionRepository) GetById(id int) (poll entity.ClientPollQuestion, err error) {
	command := `SELECT 
       				   id, 
					   client_poll_id, 
					   multiple_choice, 
					   question_text
      		 FROM client_poll_question
						 WHERE  id = $1
						 `
	err = cpqr.Conn.Get(&poll, command, id)
	return poll, err
}

func (cpqr ClientPollQuestionRepository)GetAllByIdPoll(pollID int) (questions []basic.PollQuestion, err error){
	command := `SELECT 
       					id, 
       					multiple_choice,
						question_text
				FROM client_poll_question
				WHERE client_poll_id = $1`

	err = cpqr.Conn.Select(&questions, command, pollID)

	if err != nil {
		return questions, cpqr.Error(err)
	}
	return questions, cpqr.Error(err)
}


func (cpqr ClientPollQuestionRepository) GetAggreagateById(id int) (question aggregate.PollQuestionAggregate, err error) {
	command := `SELECT 
       				   id, 
					   multiple_choice, 
					   question_text
      		 FROM client_poll_question
						 WHERE  id = $1
						 `
	err = cpqr.Conn.Get(&question.PollQuestion, command, id)
	if err != nil{
		return question, err
	}

	fmt.Println(id)

	command = `SELECT 
       				   id, 
					   choice_text
      		 FROM client_poll_choice
						 WHERE  client_poll_question_id = $1
						 `
	err = cpqr.Conn.Select(&question.Choices, command, id)
	if err != nil{
		return question, cpqr.Error(err)
	}


	return question, cpqr.Error(err)
}

func (cpqr ClientPollQuestionRepository) ValidatePoll(questionID, pollID int) (validate bool, err error) {
	var  total int
	command := `SELECT count(*)  FROM client_poll_question 
				WHERE id = $1
				AND client_poll_id = $2`
	err = cpqr.GetConnection().Get(&total,command,questionID,pollID)
	return total > 0,err
}