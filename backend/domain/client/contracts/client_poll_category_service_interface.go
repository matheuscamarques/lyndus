package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientPollCategoryServiceInterface interface {
	Create(pollCategory entity.ClientPollCategory) (id int, err error)
	Exist(pollCategory entity.ClientPollCategory) (exist bool, err error)
}
