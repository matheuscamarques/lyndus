package contracts

import (
	entity2 "bitbucket.org/lyndus/backend/domain/bonder/entity"
	"time"

	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"
)

type AppUserBSRepositoryInterface interface {
	db.ConnectorInterface
	GetBS(bsId int) (bs entity.BS, err error)
	GetBSID(bsId int) (id int, err error)
	GetBSWeekDays(bsID int) (weekDays []entity.WeekDayBS, err error)
	GetBSDaysOff(bsID int, start, end time.Time) (days []entity.DayOff, err error)
	GetBSName(bsId int) (name string, err error)
	GetBSBMName(bsId, bmID int) (name string, err error)
	GetBSCategories(bsID int) (categories []entity2.Category, err error)
	GetBSs(fantasyName string, lat, lon float64, criteria criteria.Criteria) (bss []aggregate.BSAggregate, totalPages int, err error)
	GetBSsByCategory(categoryID int, fantasyName string, lat, lon float64, criteria criteria.Criteria) (bss []aggregate.BSAggregate, totalPages int, err error)

	GetBsWeekDays(bsID int) (weekDays []aggregate.WeekDay, err error)

	GetLocation(bsID int) (location string, err error)

	GetBsIDServiceIDWeekDays(bsID, serviceId int) (weekDays []entity.WeekDayBS, err error)
	GetBSsMap(lat, lon float64) (bss []entity.BS, err error)
	GetBSsByCategoryMap(categoryID int, lat, lon float64) (bss []entity.BS, err error)
	GetBSsWithLocalization(c criteria.Criteria, lat float64, lon float64) ([]aggregate.BSAggregate, error)
	GetBSsByCategoryWithLocalization(id int, c criteria.Criteria, lat float64, lon float64) ([]aggregate.BSAggregate, error)
}
