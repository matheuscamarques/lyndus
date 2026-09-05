package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
)

type AppUserCreateService struct {
	contracts.AppUserCreateServiceInterface
	repo contracts.AppUserCreateRepositoryInterface
}

var CreateService contracts.AppUserCreateServiceInterface

func NewAppUserCreateService() *AppUserCreateService {
	return &AppUserCreateService{
		repo: repository.NewAppUserCreateRepository(),
	}
}

func (aucr AppUserCreateService) CreatePerson(auc entity.AppUserCreate) (id int, err error) {
	return aucr.repo.CreatePerson(auc)
}
func (aucr AppUserCreateService) UpdatePerson(auc entity.AppUserCreate) error {
	return aucr.repo.UpdatePerson(auc)
}
func (aucr AppUserCreateService) GetPersonIDByCPF(a *entity.AppUserCreate) (n int, err error) {
	return aucr.repo.GetPersonIDByCPF(a)
}
func (aucr AppUserCreateService) GetAppUserIDByCPF(a *entity.AppUserCreate) (n int, err error) {
	return aucr.repo.GetAppUserIDByCPF(a)
}
