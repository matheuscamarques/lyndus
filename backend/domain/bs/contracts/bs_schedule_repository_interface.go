package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type BSScheduleRepositoryInterface interface {
	CreateScheduling(scheduling entity.Scheduling) (id int, err error)
	GetBSScheduling(bsId int, date time.Time) (scheduling []entity.Scheduling, err error)
	ChangeStatusSchedule(bsID, scheduleID, status int) (err error)
	GetBSScheduledTime(bsId, bmID int, startTime, endTime types.DateTime) (ids []int, err error)
	GetBSSchedule(bsId, scheduleID int) (scheduling entity.Scheduling, err error)
	SetOrderIDSchedule(bsID, scheduleID, orderID int) (err error)
}
