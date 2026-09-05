package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BSAppUserService struct {
	repo contracts.BSAppUserRepositoryInterface
}

var AppUserService contracts.BSAppUserServiceInterface


func NewBSAPPuserService() *BSAppUserService {
	return &BSAppUserService{repo: repository.NewBSAppUserRepository()}
}

func (baus BSAppUserService) GetAppUserByCPF(cpf types.CPF) (appUser entity.AppUser, err error) {
	return baus.repo.GetAppUserByCPF(cpf)
}

func (baus BSAppUserService) GetAppUserIDByCPF(cpf *types.CPF) (id int, err error) {
	return baus.repo.GetAppUserIDByCPF(cpf)
}
