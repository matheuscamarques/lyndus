package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
	"errors"
)

type AppUserUserService struct {
	repo contracts.AppUserUserRepositoryInterface
}

var UserUserService contracts.AppUserUserServiceInterface

func NewAppUserUserService() *AppUserUserService {
	return &AppUserUserService{
		repo: repository.NewAppUserUserRepository(),
	}
}
func (aur AppUserUserService) UpdateAppUser(appUser entity.AppUser) error {
	return aur.repo.UpdateAppUser(appUser)
}
func (aur AppUserUserService) CreateAppUser(appUser entity.AppUser) (id int, err error) {
	return aur.repo.CreateAppUser(appUser)
}
func (aur AppUserUserService) GetAppUserByID(id int) (appUser entity.AppUserCreate, err error) {
	return aur.repo.GetAppUserByID(id)
}
func (aur AppUserUserService) GetAppUserFullByID(id int) (appUser entity.AppUserCreate, err error) {
	return aur.repo.GetAppUserFullByID(id)
}
func (aur AppUserUserService) UpdateAppUserProfile(appUser entity.AppUserCreate) (err error) {

	appUser.PersonID, err = aur.repo.GetAppUserPersonIDByID(appUser.ID)
	if err != nil {
		return
	}
	if appUser.PersonID == 0 {
		err = errors.New("not find bs_person")
	}
	err = aur.repo.UpdateAppUserPerson(appUser)
	if err != nil {
		return
	}
	err = aur.repo.UpdateAppUserEmail(appUser)
	if err != nil {
		return
	}
	return
}
