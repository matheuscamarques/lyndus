package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSWeekDayServiceInterface interface {
	GetWeekDays() (weekDays []entity.WeekDay, err error)
	Get(bsID int, weekDayID int) (wds []entity.WeekDayBS, err error)
}
