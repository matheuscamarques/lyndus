package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/types"
)

type AppUserBSPersonRepositoryInterface interface {
	GetBSPersonByCPF(bsID int, cpf types.CPF) (bm entity.BSPerson, err error)
	GetBSPersonByAppUserID(bsID, appUserID int) (bm entity.BSPerson, err error)
	CreatePerson(person entity.BSPerson) (id int, err error)
	UpdatePersonAppUserID(bsID, personID, appUserID int) error
}
