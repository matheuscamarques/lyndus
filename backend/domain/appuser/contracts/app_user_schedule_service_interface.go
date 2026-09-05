package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/global/basic"
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type AppUserScheduleServiceInterface interface {
	CreateScheduling(scheduling entity.Scheduling) (id int, err error)
	GetBSScheduling(bsId int, date time.Time) (scheduling []entity.Scheduling, err error)
	ChangeStatusSchedule(bsID, scheduleID, status int) (err error)
	GetBSScheduledTime(bsId, bmID int, startTime, endTime types.DateTime) (ids []int, err error)
	GetBSBMScheduling(bsId, bmID int, start, end time.Time) (scheduling []entity.SchedulingBAsic, err error)

	GetSchedulingsByAppUserID(appUserID, activePage, itemsPerPage int, date types.Date) (scheduling []basic.AppScheduling, totalPages int, err error)
	GetSchedulingAppUser(appUserID, scheduleID int) (scheduling entity.Scheduling, err error)
}
