package contracts

import "bitbucket.org/lyndus/backend/domain/appuser/entity"

type AppUserUserServiceInterface interface {
	UpdateAppUser(appUser entity.AppUser) error
	UpdateAppUserProfile(appUser entity.AppUserCreate) error
	CreateAppUser(appUser entity.AppUser) (id int, err error)
	GetAppUserByID(id int) (appUser entity.AppUserCreate, err error)
	GetAppUserFullByID(id int) (appUser entity.AppUserCreate, err error)

}
