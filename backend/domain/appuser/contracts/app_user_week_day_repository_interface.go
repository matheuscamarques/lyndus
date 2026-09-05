package contracts

import "bitbucket.org/lyndus/backend/domain/appuser/entity"

type AppUserWeekDayRepositoryInterface interface {
	// GetWeekDays this method return week of days
	GetWeekDays() (weekDays []entity.WeekDay, err error)
	Get(bsID, weekDayID int) (wds []entity.WeekDayBS, err error)
}
