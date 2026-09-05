package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/global/basic"
	"bitbucket.org/lyndus/backend/global/contracts"
	"bitbucket.org/lyndus/backend/infra/criteria"
)

type ClientPollRepositoryInterface interface {
	contracts.RepositoryInterface
	Create(poll entity.ClientPoll) (id int, err error)
	Update(poll entity.ClientPoll) (err error)
	GetById(id int) (poll aggregate.PollAgregate, err error)
	GetAllByIdClient(clientID int, criteria criteria.Criteria) (polls []basic.Poll, err error)
	ValidateClient(idPoll, idClient int) (bool, error)
}
