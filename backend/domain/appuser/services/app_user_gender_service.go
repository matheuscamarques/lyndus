package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
)

type AppUserGenderService struct {
	repo contracts.AppUserGenderRepositoryInterface
}

var UserGenderService contracts.AppUserGenderServiceInterface

func NewAppAppUserGenderService() *AppUserGenderService {
	return &AppUserGenderService{
		repo: repository.NewAppUserGenderRepository(),
	}
}

func (augs AppUserGenderService) GetGenders() (genders []entity.Gender, err error) {
	return augs.repo.GetGenders()
}

func (augs AppUserGenderService) GetCategories() (categories []entity.Category, err error) {
	return augs.repo.GetCategories()
}

func (augs AppUserGenderService)GetGenderID(genderID int) (id int, err error){
	return augs.repo.GetGenderID(genderID)
}

func (augs AppUserGenderService)GetGenderByID(genderID int) (gender entity.Gender, err error){
	return augs.repo.GetGenderByID(genderID)
}