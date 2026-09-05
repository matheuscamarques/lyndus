package services

import (
	"bitbucket.org/lyndus/backend/domain/bonder/contracts"
	"bitbucket.org/lyndus/backend/domain/bonder/repository"
	"bitbucket.org/lyndus/backend/domain/bs/services"
)

type BonderBSUserService struct {
	services.BSUserService
}

var BsUser contracts.BonderBSUserServiceInterface

func NewBonderBSUserService()  *BonderBSUserService{
	repo := repository.NewBSUserRepository()
	service := BonderBSUserService{}
	service.SetRepo(repo)
	return &service
}

func (bs BonderBSUserService)GetMasterUserAuthID(id int) (int, error){
	return bs.Repo.GetMasterUserAuthID(id)
}