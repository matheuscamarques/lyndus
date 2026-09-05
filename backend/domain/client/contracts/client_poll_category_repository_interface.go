package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientPollCategoryRepositoryInterface interface {
	Create(entity.ClientPollCategory) (id int, err error)
	Update(entity.ClientPollCategory) (err error)
	GetById(id int) (pollCategory entity.ClientPollCategory, err error)
	Exist(pollCategory entity.ClientPollCategory) (exist bool, err error)
}
