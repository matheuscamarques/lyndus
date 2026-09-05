package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
)

type AppUserServiceRepositoryInterface interface {
	GetServicesByCategory(bsID, serviceCategoryID int) (services []entity.Service, err error)
	GetServiceCategory() (serviceCategories []entity.ServiceCategory, err error)
}
