package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/types"
)

type AppUserServiceInterface interface {
	GetAppUserByCPF(cpf types.CPF, page int, pag int) (response criteria.AppUserResponse, err error)
	GetAppUserByName(name string, page int, pag int) (response criteria.AppUserResponse, err error)

	//GetAppUsers(page int, pag int) (response criteria.AppUserResponse, err error)
	GetAllAppUser(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (appUsers []entity.AppUserPerson, totalItems, totalPages int, err error)

	UpdateLyndusBox(bonus entity.Bonus) error

	GetAppUserByID(appUserID int) (response entity.AppUser, err error)

	//UpdateBonus(bonus entity.Bonus) error

	//GetAllAppUser(page int, pag int, search string) (response criteria.AppUserResponse, err error)
}
