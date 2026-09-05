package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BSPersonRepositoryInterface interface {
	FindPerson(bsID int, cpf *types.CPF, phone string) (person entity.Person, err error)
	GetPersonID(bsID, personID int) (person entity.Person, err error)
	GetPerson(bsID, personID int) (person entity.Person, err error)
	GetPersons(bsID int) (persons []entity.Person, err error)
	GetPersonsByCPF(bsID int, cpf types.CPF) (persons []entity.Person, err error)
	GetPersonsByPhone(bsID int, phone string) (persons []entity.Person, err error)
	GetPersonsByName(bsID int, name string) (persons []entity.Person, err error)
	CreatePerson(person entity.Person) (id int, err error)
	UpdatePerson(person entity.Person) error
	UpdatePersonCPF(bsID, personID int, cpf types.CPF) error
	UpdatePersonCPFAndAppUser(bsID, personID int, cpf types.CPF, appUserID int) error
}
