package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientPollChoiceRepositoryInterface interface {
	Create(pollChoice entity.ClientPollChoice)(id int,err error)
	Update(pollChoice entity.ClientPollChoice)(err error)
	GetById(id int)(pollChoice entity.ClientPollChoice,err error)
}
