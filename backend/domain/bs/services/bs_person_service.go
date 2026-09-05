package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BSPersonService struct {
	repo contracts.BSPersonRepositoryInterface
}
var PersonService contracts.BSPersonServiceInterface

func NewBSPersonService() *BSPersonService {
	return &BSPersonService{
		repo: repository.NewBSPersonRepository(),
	}
}

func (bps BSPersonService) FindPerson(bsID int, cpf *types.CPF, phone string) (person entity.Person, err error) {
	return bps.repo.FindPerson(bsID, cpf, phone)
}
func (bps BSPersonService) GetPersonID(bsID, personID int) (person entity.Person, err error) {
	return bps.repo.GetPersonID(bsID, personID)
}
func (bps BSPersonService) GetPerson(bsID, personID int) (person entity.Person, err error) {
	return bps.repo.GetPerson(bsID, personID)
}
func (bps BSPersonService) GetPersons(bsID int) (persons []entity.Person, err error) {
	return bps.repo.GetPersons(bsID)
}
func (bps BSPersonService) GetPersonsByCPF(bsID int, cpf types.CPF) (persons []entity.Person, err error) {
	return bps.repo.GetPersonsByCPF(bsID, cpf)
}
func (bps BSPersonService) GetPersonsByPhone(bsID int, phone string) (persons []entity.Person, err error) {
	return bps.repo.GetPersonsByPhone(bsID, phone)
}
func (bps BSPersonService) GetPersonsByName(bsID int, name string) (persons []entity.Person, err error) {
	return bps.repo.GetPersonsByName(bsID, name)
}
func (bps BSPersonService) CreatePerson(person entity.Person) (id int, err error) {
	return bps.repo.CreatePerson(person)
}
func (bps BSPersonService) UpdatePerson(person entity.Person) error {
	return bps.repo.UpdatePerson(person)
}
func (bps BSPersonService) UpdatePersonCPF(bsID, personID int, cpf types.CPF) error {
	return bps.repo.UpdatePersonCPF(bsID, personID, cpf)
}
func (bps BSPersonService) UpdatePersonCPFAndAppUser(bsID, personID int, cpf types.CPF, appUserID int) error {
	return bps.repo.UpdatePersonCPFAndAppUser(bsID, personID, cpf, appUserID)
}
