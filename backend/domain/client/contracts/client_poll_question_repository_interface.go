package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/global/basic"
	"bitbucket.org/lyndus/backend/infra/db"
	"github.com/jmoiron/sqlx"
)

type ClientPollQuestionRepositoryInterface interface {
	Create(poll entity.ClientPollQuestion) (id int, err error)
	Update(poll entity.ClientPollQuestion) (err error)
	GetById(id int) (question entity.ClientPollQuestion, err error)
	GetAggreagateById(id int)(question aggregate.PollQuestionAggregate, err error)
	GetAllByIdPoll(pollID int) (questions []basic.PollQuestion, err error)
	GetConnector() db.Connector
	GetTable() string
	GetConnection() *sqlx.DB
	ValidatePoll(questionID, pollID int) (bool, error)
}



