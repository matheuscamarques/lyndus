package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientPersonRepositoryInterface interface {
	CreatePersonBasic(person *entity.PersonBasic) error
	GetPersonIDByCPF(person *entity.PersonBasic) error
}
