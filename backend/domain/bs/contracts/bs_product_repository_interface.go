package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type BSProductRepositoryInterface interface {
	db.ConnectorInterface
	GetCategories() (categories []entity.ProductCategory, err error)
	GetPackages() (packages []entity.ProductPackage, err error)
	GetUnits() (units []entity.ProductUnit, err error)
	GetClasses() (classes []entity.ProductClass, err error)
	Update(product entity.Product) error
	Create(product entity.Product) (int, error)
	GetByID(bsID, productID int) (product entity.Product, err error)
	GetAllProducts(bsID int, c criteria.Criteria) ([]entity.Product, error)
	CountProducts(bsID int) (n int, err error)
	SetConn(conn *sqlx.DB)
	StockSave(bsId, createdBy, id int, quantity int, costPrice, salesPrice, profitMargin decimal.Decimal) error
	GetProductsSimple(bsID int) (products []entity.ProductSimple, err error)
	GetProductSimple(bsID, productID int) (product entity.ProductSimple, err error)

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
