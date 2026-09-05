package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSCompanyRepositoryInterface interface {
	GetCompanyByBsID(bsID int) (company entity.Company, err error)
	UpdateCompanyByBS(company entity.Company) error
	GetBSWeekDays(bsID int) (weekDays []entity.BMWeekDay, err error)
	DeleteBSAllWeekDays(bsID int) error
	AddBSWeekDay(bsID int, day entity.BMWeekDay) error
	GetLocation(bsID int) (location string, err error)
}
