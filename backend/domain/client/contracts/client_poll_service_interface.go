package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/criteria"
)

type ClientPollServiceInterface interface {
	Create(poll entity.ClientPoll)(id int, err error)
	GetById(pollID int)(poll aggregate.PollAgregate,err error)
	GetAll(clientID ,ActivePage,ItemsPerPage int)(response criteria.CPollResponse,err error)
	Update(poll entity.ClientPoll)(err error)
	ValidateClient(pollID , clientID int) (bool, error)
}
