package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientPollResultsRepositoryInterface interface {
	Create(pollResults entity.ClientPollResults) (id int, err error)
	Update(pollResults entity.ClientPollResults) (err error)
	GetById(id int) (pollResults entity.ClientPollResults, err error)
}
