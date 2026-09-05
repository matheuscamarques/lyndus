package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/global/contracts"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BonderClientRepositoryInterface interface {
	contracts.RepositoryInterface
	Create(user entity.Client) (int, error)
	SearchByCNPJ(cnpj types.CNPJ) (id int, err error)
	GetAll(activePage, itemsPerPage int, search, orderby string, sortDesc, active bool) (companies []entity.Company, totalItems, totalPages int, err error)
	GetById(id int) (entity.Client, error)
	GetCompanyId(id int) (int, error)

	UpdateLyndusActive(bsId int, active bool) (err error)
	Detete(bsId int) (err error)
}
