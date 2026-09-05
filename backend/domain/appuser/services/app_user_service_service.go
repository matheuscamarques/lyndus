package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
)

type AppUserServiceService struct {
	repo contracts.AppUserServiceRepositoryInterface
}

var AppUserService contracts.AppUserServiceInterface

func NewAppUserServiceService() *AppUserServiceService {
	return &AppUserServiceService{
		repo: repository.NewAppUserServiceRepository(),
	}
}

func (ass AppUserServiceService) GetServiceCategory() (serviceCategories []entity.ServiceCategory, err error){
	return ass.repo.GetServiceCategory()
}

func (ass AppUserServiceService) GetServicesByCategory(bsID, serviceCategoryID int) (services []entity.Service, err error){
	return ass.repo.GetServicesByCategory(bsID, serviceCategoryID)
}
