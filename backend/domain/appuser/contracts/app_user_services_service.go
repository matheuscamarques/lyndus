package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
)

type AppUserServiceInterface interface {

	//FilterServicesCategoriesByDate(servicesCategories []entity.ServiceCategory,hour types.TimeHHMM) ([]entity.ServiceCategory, error)
	//FilterServicesCategoriesByHour(servicesCategories []entity.ServiceCategory,hour types.TimeHHMM) ([]entity.ServiceCategory, error)
	GetServicesByCategory(bsID, serviceCategoryID int) (services []entity.Service, err error)
	GetServiceCategory() (serviceCategories []entity.ServiceCategory, err error)
}
