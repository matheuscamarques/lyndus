package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
)

type ClientAppUserService struct {
	contracts.ClientAppUserServiceInterface
	repo contracts.ClientAppUserRepositoryInterface
}

var AppUserService contracts.ClientAppUserServiceInterface

func NewClientAppUserService() *ClientAppUserService {
	return &ClientAppUserService{
		repo: repository.NewClientAppUserRepository(),
	}
}

func (cpus ClientAppUserService) CreateAppUserBasic(user *entity.AppUser) error {
	return cpus.repo.CreateAppUserBasic(user)
}

func (cpus ClientAppUserService) GetAppUserIDByCPF(user *entity.AppUser) error {
	return cpus.repo.GetAppUserIDByCPF(user)
}
