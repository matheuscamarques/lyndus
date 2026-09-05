package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BSAppUserRepositoryInterface interface {
	GetAppUserByCPF(cpf types.CPF) (appUser entity.AppUser, err error)
	GetAppUserIDByCPF(cpf *types.CPF) (id int, err error)
}
