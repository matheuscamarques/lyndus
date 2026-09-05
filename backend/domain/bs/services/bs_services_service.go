package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BSServicesService struct {
	repo contracts.BSServiceRepositoryInterface
}



var BS contracts.BSServiceServiceInterface

func NewBSServicesService() *BSServicesService {
	return &BSServicesService{
		repo: repository.NewBSServiceRepository(),
	}
}

func (bss BSServicesService) GetServiceID(bsID, serviceID int) (n int, err error) {
	return bss.repo.GetServiceID(bsID, serviceID)
}
func (bss BSServicesService) GetService(bsID, serviceID int) (service entity.Service, err error) {
	return bss.repo.GetService(bsID, serviceID)
}
func (bss BSServicesService) GetBSServices(bsID int) (services []entity.Service, err error) {
	return bss.repo.GetBSServices(bsID)
}
func (bss BSServicesService) UpdateService(service entity.Service) error {
	return bss.repo.UpdateService(service)
}
func (bss BSServicesService)GetBSServiceCategories() (serviceCategories []entity.ServiceCategory, err error) {
	return bss.repo.GetBSServiceCategories()
}
func (bss BSServicesService) GetBSServicesByName(bsID int, name string) (services []entity.Service, err error) {
	return bss.repo.GetBSServicesByName(bsID, name)
}
func (bss BSServicesService) CreateService(service entity.Service) (id int, err error) {
	return bss.repo.CreateService(service)
}
func (bss BSServicesService) AddBMService(bmID int ,service entity.Service) error {
	return bss.repo.AddBMService(bmID, service)
}
func (bss BSServicesService) DeleteBMService(bmID, serviceID int) error {
	return bss.repo.DeleteBMService(bmID, serviceID)
}

func (bss BSServicesService) GetBankDetails(id int) (entity.BSBankDetails, error) {
	return bss.repo.GetBankDetails(id)
}

func (bss BSServicesService) UpdateBankDetails(bank entity.BSBankDetails) error {
	return bss.repo.UpdateBankDetails(bank)
}

func (bss BSServicesService) CreateBankDetails(bank entity.BSBankDetails)  error {
	return bss.repo.CreateBankDetails(bank)
}