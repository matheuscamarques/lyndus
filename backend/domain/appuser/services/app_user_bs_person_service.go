package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
	"bitbucket.org/lyndus/backend/infra/types"
)

type AppUserBSPersonService struct {
	repo contracts.AppUserBSPersonRepositoryInterface
}

var PersonService contracts.AppUserBSPersonServiceInterface

func NewAppUserBSPersonService() *AppUserBSPersonService {
	return &AppUserBSPersonService{
		repo: repository.NewAppUserBSPersonRepository(),
	}
}

func (aup AppUserBSPersonService) GetBSPersonByCPF(bsID int, cpf types.CPF) (bm entity.BSPerson, err error) {
	return aup.repo.GetBSPersonByCPF(bsID, cpf)
}
func (aup AppUserBSPersonService) GetBSPersonByAppUserID(bsID, appUserID int) (bm entity.BSPerson, err error) {
	return aup.repo.GetBSPersonByAppUserID(bsID, appUserID)
}

func (aup AppUserBSPersonService) CreatePerson(person entity.BSPerson) (id int, err error) {
	return aup.repo.CreatePerson(person)
}

func (aup AppUserBSPersonService) UpdatePersonAppUserID(bsID, personID, appUserID int) error {
	return aup.repo.UpdatePersonAppUserID(bsID, personID, appUserID)
}