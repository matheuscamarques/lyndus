package contracts

import (
	"bitbucket.org/lyndus/backend/global/aggregate"
	"time"

	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/criteria"
)

type AppUserBSServiceInterface interface {
	GetBS(bsId int) (bs entity.BS, err error)
	GetBSID(bsId int) (id int, err error)
	GetBSWeekDays(bsID int) (weekDays []entity.WeekDayBS, err error)
	GetBSDaysOff(bsID int, start, end time.Time) (days []entity.DayOff, err error)

	GetBSName(bsId int) (name string, err error)
	GetBSBMName(bsId, bmID int) (name string, err error)
	GetBSs(page, itemsPerPage int, fantasyName string, lat, lon float64) (response criteria.BSResponse, err error)
	GetBSsByCategory(categoryID, page, itemsPerPage int, fantasyName string, lat, lon float64) (response criteria.BSResponse, err error)

	GetBsIDServiceIDWeekDays(bsID, serviceId int) (weekDays []entity.WeekDayBS, err error)
	GetBSsMap(lat, lon float64) (bss []entity.BS, err error)
	GetBSsByCategoryMap(categoryID int, lat, lon float64) (bss []entity.BS, err error)

	GetBsWeekDays(bsID int) (weekDays []aggregate.WeekDay, err error)

	GetLocation(bsID int) (location string, err error)

	GetBSsWithLocalization(lat, lon float64, page, itemsPerPage int) (response criteria.BSResponse, err error)
	GetBSsByCategoryLocalization(categoryID int, lat, lon float64, page, itemsPerPage int) (response criteria.BSResponse, err error)
}
