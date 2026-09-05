package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BonderClientServiceInterface interface {
	Register(entity.Client) (int, error)
	SearchByCNPJ(cnpj types.CNPJ) (int, error)
	GetAll(activePage, itemsPerPage int, search, orderby string, sortDesc, active bool) (companies []entity.Company, totalItems, totalPages int, err error)
	GetById(id int) (entity.Client, error)
	GetCompanyId(id int) (int, error)

	UpdateLyndusActive(clientId int, active bool) (err error)
	Detete(clientId int) (err error)
}
