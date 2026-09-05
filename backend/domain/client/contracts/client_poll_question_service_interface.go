package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/global/basic"
)

type ClientPollQuestionServiceInterface interface {
	Create(pollQuestion entity.ClientPollQuestion)(int , error)
	Replicate(pollQuestion entity.ClientPollQuestion)(int , error)
	Update(pollQuestion entity.ClientPollQuestion) (err error)
	GetAll(pollID int)(questions []basic.PollQuestion, err error)
	GetById(questionID int)(question aggregate.PollQuestionAggregate,err error)
	ValidatePoll(questionID , pollID int) (bool, error)
}
