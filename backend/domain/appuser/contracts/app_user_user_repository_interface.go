package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
)

type AppUserUserRepositoryInterface interface {
	UpdateAppUserEmail(appUser entity.AppUserCreate) error
	UpdateAppUserPerson(appUser entity.AppUserCreate) error
	UpdateAppUser(appUser entity.AppUser) error
	CreateAppUser(appUser entity.AppUser) (id int, err error)
	GetAppUserByID(id int) (appUser entity.AppUserCreate, err error)
	GetAppUserFullByID(id int) (appUser entity.AppUserCreate, err error)
	GetAppUserPersonIDByID(id int) (personID int, err error)
}
