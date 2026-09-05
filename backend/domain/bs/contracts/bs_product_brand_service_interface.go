package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/criteria"
)

type BSProductBrandServiceInterface interface {
	Create(entity.ProductBrand) (int, error)
	GetByID(brandID int) (entity.ProductBrand, error)
	GetAll(page, itemsPerPag int)(criteria.BSProductBrand, error)
	GetByName(name string)(brands entity.ProductBrand, err error)
	SearchByName(name string) (brands []entity.ProductBrand, err error)
}
