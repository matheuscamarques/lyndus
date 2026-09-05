package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/global/contracts"
	"bitbucket.org/lyndus/backend/infra/criteria"
)

type BSProductBrandRepositoryInterface interface {
	contracts.RepositoryInterface
	Create(brand entity.ProductBrand) (id int, err error)
	Update(brand entity.ProductBrand) error
	GetByID(brandID int) (brand entity.ProductBrand, err error)
	GetAll(criteria criteria.Criteria) ([]aggregate.BSProductBrandAggregate, error)
	GetByName(name string) (brands entity.ProductBrand, err error)
	SearchByName(name string) (brands []entity.ProductBrand, err error)

	//FindPerson(bsID int, cpf types.CPF, phone string) (bs_person entity.Person, err error)
	//GetPersonID(bsID, personID int) (bs_person entity.Person, err error)
	//GetPerson(bsID, personID int) (bs_person entity.Person, err error)
	//GetPersons(bsID int) (persons []entity.Person, err error)
	//GetPersonsByCPF(bsID int, cpf types.CPF) (persons []entity.Person, err error)
	//GetPersonsByPhone(bsID int, phone string) (persons []entity.Person, err error)
	//GetPersonsByName(bsID int, name string) (persons []entity.Person, err error)
	//CreatePerson(bs_person entity.Person) (id int, err error)
	//UpdatePerson(bs_person entity.Person) error
}
