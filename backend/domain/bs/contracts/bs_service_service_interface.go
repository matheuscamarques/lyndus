package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSServiceServiceInterface interface {
	GetServiceID(bsID, serviceID int) (n int, err error)
	GetService(bsID, serviceID int) (service entity.Service, err error)
	GetBSServices(bsID int) (services []entity.Service, err error)
	UpdateService(service entity.Service) error
	GetBSServiceCategories() (serviceCategories []entity.ServiceCategory, err error)
	GetBSServicesByName(bsID int, name string) (services []entity.Service, err error)
	CreateService(service entity.Service) (id int, err error)
	AddBMService(bmID int,service entity.Service) error
	DeleteBMService(bmID, serviceID int) error
	GetBankDetails(id int) (entity.BSBankDetails, error)
	UpdateBankDetails(bank entity.BSBankDetails) error
	CreateBankDetails(bank entity.BSBankDetails) error
}
