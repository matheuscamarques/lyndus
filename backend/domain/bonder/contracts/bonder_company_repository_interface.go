package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BonderCompanyRepositoryInterface interface {
	GetIdByCNPJ(cnpj types.CNPJ) (int, error)
	Create(company *entity.BS) (int, error)
	Update(company *entity.BS) error
}
