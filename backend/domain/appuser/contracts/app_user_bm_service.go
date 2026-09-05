package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"time"
)

type AppUserBMInterface interface {
	GetBMsByServiceID(bsID, serviceID int) (bms []entity.BM, err error)
	GetBMsBasicByServiceID(bsID, serviceID int) (bms []entity.BM, err error)
	GetBMWeekDays(bmID int) (weekDays []entity.WeekDayBS, err error)
	GetBMWeekDaysWithName(bmID int) (weekDays []aggregate.WeekDay, err error)
	GetBMDaysOff(bsID, bmID int, start, end time.Time) (days []entity.DayOff, err error)
	GetBMByID(bsID, bmID int) (bm entity.BM, err error)
	GetBMByIDAndDuration(bsID, bmID, serviceID int) (bm entity.BM, err error)
	GetBMsFullByWeekDay(bsID, serviceID, weekDayID int) (bm []entity.BM, err error)
	GetBMFullByWeekDay(bsID, serviceID, bmID,  weekDayID int) (bm entity.BM, err error)
}
