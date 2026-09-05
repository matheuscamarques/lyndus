package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BSBMServiceInterface interface {
	GetBMService(bmID, serviceID int) (entity.Service, error)
	GetBmID(bsID, bmID, bmStatusID int) (n int, err error)
	GetBM(bsID, bmID, bmStatusID int) (bm entity.BM, err error)
	GetBmIDByCPF(bsID int, cpf types.CPF) (id int, err error)
	CreateBM(bm entity.BM) (id int, err error)
	UpdateBM(bm entity.BM) error
	ChangeStatus(bsID, bmID, bmStatusID int) error
	AddBMService(bmID int, service entity.Service) error
	DeleteBMService(bmID, serviceID int) error
	AddBMWeekDay(bmID int, day entity.BMWeekDay) error
	DeleteBMAllWeekDays(bmID int) error
	GetBMS(bsID, page, itemsPerPage int, search string) (bms []entity.BM, totalItems, totalPages int, err error)
	GetBmServicesIDs(bmID int) (services []int, err error)
	GetBmServices(bsID, bmID int) (services []entity.Service, err error)
	GetBMServices(bmID int) (services []entity.Service, err error)
	GetBMWeekDays(bmID int) (weekDays []entity.BMWeekDay, err error)
	GetBMsByWeekDay(bsID, weekDay int) (weekDays []entity.BMBasic, err error)

	UpdateBMService(bmID int, service entity.Service) error
}
