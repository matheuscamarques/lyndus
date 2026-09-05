package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
)

type ClientPersonService struct {
	contracts.ClientPersonServiceInterface
	repo contracts.ClientPersonRepositoryInterface
}

var PersonService contracts.ClientPersonServiceInterface

func NewClientPersonService() *ClientPersonService {
	return &ClientPersonService{
		repo: repository.NewClientPersonRepository(),
	}
}

func (cps ClientPersonService) CreatePersonBasic(person *entity.PersonBasic) error {
	return cps.repo.CreatePersonBasic(person)
}
func (cps ClientPersonService) GetPersonIDByCPF(person *entity.PersonBasic) error {
	return cps.repo.GetPersonIDByCPF(person)
}
