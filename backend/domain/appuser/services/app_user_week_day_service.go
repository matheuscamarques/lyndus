package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
)

type AppUserWeekDayService struct {
	repo contracts.AppUserWeekDayRepositoryInterface
}

var WeekDayService contracts.AppUserWeekDayServiceInterface

func NewAppUserWeekDayService() *AppUserWeekDayService {

	return &AppUserWeekDayService{
		repo: repository.NewAppUserWeekDayRepository(),
	}
}

func (wdrs AppUserWeekDayService) GetWeekDays() (weekDays []entity.WeekDay, err error) {
	return wdrs.repo.GetWeekDays()
}

func (wdrs AppUserWeekDayService) Get(bsID int, weekDayID int) (wds []entity.WeekDayBS, err error){
	return wdrs.repo.Get(bsID,weekDayID)
}