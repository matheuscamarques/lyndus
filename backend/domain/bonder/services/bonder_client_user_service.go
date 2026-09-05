package services

import (
	"bitbucket.org/lyndus/backend/domain/bonder/contracts"
	"bitbucket.org/lyndus/backend/domain/bonder/repository"
	"bitbucket.org/lyndus/backend/domain/client/services"
)

type BonderClientUserService struct {
	services.ClientUserService
}

var ClientUser contracts.BonderClientUserServiceInterface

func NewBonderClientUserService()  *BonderClientUserService {
	repo    := repository.NewBonderClientUserRepository()
	service := BonderClientUserService{}
	service.SetRepo(&repo)
	return &service
}

func (cu BonderClientUserService)GetMasterUserAuthID(ClientID int) (int, error){
	return cu.Repo.GetMasterUserAuthID(ClientID)
}