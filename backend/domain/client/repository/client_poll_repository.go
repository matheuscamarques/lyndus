package repository

import (
	"fmt"

	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/global/basic"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"
	"github.com/jmoiron/sqlx"
)

type ClientPollRepository struct {
	db.Connector
}

func NewClientPollRepository() (repo ClientPollRepository) {
	repo = ClientPollRepository{}
	repo.UsePostgres()
	repo.Table = "client_poll"
	return repo
}

func (cpr ClientPollRepository) Create(poll entity.ClientPoll) (id int, err error) {
	command := `INSERT INTO 
							client_poll(
										client_id,
							            start_time, 
							            end_time, 
							            all_employee, 
							            created_by, 
							            title, 
							            status
											 	)
				VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`
	stmt, err := cpr.Conn.Prepare(command)
	if err != nil {
		return id, cpr.Error(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		poll.ClientID,
		poll.StartTime,
		poll.EndTime,
		poll.AllEmployee,
		poll.CreatedBy,
		poll.Title,
		poll.Status,
	).Scan(&id)
	return id, cpr.Error(err)
}
func (cpr ClientPollRepository) Update(poll entity.ClientPoll) (err error) {
	command := `UPDATE client_poll 
								 SET start_time = $3,
								     end_time = $4,
								     all_employee = $5,
								     title = $6,
								     status = $7
								 WHERE client_id=$1
								 AND id=$2`
	stmt, err := cpr.Conn.Prepare(command)
	if err != nil {
		return cpr.Error(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		poll.ClientID,
		poll.ID,
		poll.StartTime,
		poll.EndTime,
		poll.AllEmployee,
		poll.Title,
		poll.Status,
	)
	return cpr.Error(err)
}
func (cpr ClientPollRepository) GetById(id int) (poll aggregate.PollAgregate, err error) {
	basePoll := basic.Poll{}
	command := `SELECT 
       					id, 
       					start_time, 
       					end_time, 
       					all_employee, 	
       					title, 
       					status
      		 FROM client_poll
						 WHERE  id = $1
						 `
	err = cpr.Conn.Get(&basePoll, command, id)

	if err != nil {
		return poll, cpr.Error(err)
	}
	poll.Poll = basePoll

	var questions []aggregate.PollQuestionAggregate
	command = `SELECT 
       					id,
       					multiple_choice,
       					question_text 
				FROM client_poll_question
				WHERE client_poll_id = $1`

	err = cpr.Conn.Select(&questions, command, id)

	if err != nil {
		return poll, cpr.Error(err)
	}

	poll.Questions = questions
	fmt.Printf("%+v", poll.Questions)
	for i := range poll.Questions {
		var choices []basic.PollChoice
		command = `SELECT 
       					id,
       					choice_text
				FROM client_poll_choice
				WHERE client_poll_question_id = $1`

		err = cpr.Conn.Select(&choices, command, poll.Questions[i].ID)

		if err != nil {
			return poll, cpr.Error(err)
		}
		poll.Questions[i].Choices = choices
	}

	var category []int
	command = `SELECT 
       					client_category_id
				FROM client_poll_category
				WHERE client_poll_id = $1`

	err = cpr.Conn.Select(&category, command, id)

	if err != nil {
		return poll, cpr.Error(err)
	}
	poll.Category = category
	return poll, cpr.Error(err)
}

func (cpr ClientPollRepository) GetAllByIdClient(clientID int, criteria criteria.Criteria) (polls []basic.Poll, err error) {
	command := criteria.Query(`SELECT 
       					id, 
       					start_time, 
       					end_time, 
       					all_employee, 
       					title, 
       					status
				FROM client_poll
				WHERE client_id = $1`)

	err = cpr.Conn.Select(&polls, command, clientID)

	if err != nil {
		return polls, cpr.Error(err)
	}
	return polls, cpr.Error(err)
}

func (cpr ClientPollRepository) ValidateClient(idPoll, idClient int) (validate bool, err error) {
	var total int
	command := `SELECT count(*)  FROM client_poll 
				WHERE id = $1
				AND client_id = $2`
	err = cpr.GetConnection().Get(&total, command, idPoll, idClient)
	return total > 0, err
}

func (cpr ClientPollRepository) GetConnector() db.Connector {
	return cpr.Connector
}
func (cpr ClientPollRepository) GetTable() string {
	return cpr.Table
}
func (cpr ClientPollRepository) GetConnection() *sqlx.DB {
	return cpr.Conn
}
