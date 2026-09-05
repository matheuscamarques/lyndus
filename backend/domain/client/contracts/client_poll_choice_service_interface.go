package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientPollChoiceServiceInterface interface {
	Create(pollChoice entity.ClientPollChoice)(int , error)
	Replicate(pollChoice entity.ClientPollChoice)(int , error)
	Update(pollChoice entity.ClientPollChoice) error
}
