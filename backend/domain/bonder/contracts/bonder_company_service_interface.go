package contracts

import "bitbucket.org/lyndus/backend/domain/bonder/entity"

type BonderCompanyServiceInterface interface {
	Register(company *entity.BS) (int, error)
	Update(company *entity.BS) error
}
