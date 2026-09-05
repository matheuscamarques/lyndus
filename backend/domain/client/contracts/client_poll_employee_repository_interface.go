package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientPollEmployeeRepositoryInterface interface {
	Create(poll entity.ClientPollEmployee) (id int, err error)
	Update(poll entity.ClientPollEmployee) (err error)
	GetById(id int) (poll entity.ClientPollEmployee, err error)
}
