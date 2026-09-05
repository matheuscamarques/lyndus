package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BSWeekDayService struct {
	contracts.BSWeekDayServiceInterface
	repo contracts.BSWeekDayRepositoryInterface
}

var WeekDayService contracts.BSWeekDayServiceInterface

func NewBSWeekDayService() *BSWeekDayService {
	repo := repository.NewBSWeekDayRepository()
	return &BSWeekDayService{
		repo: repo,
	}
}

func (wdrs BSWeekDayService) GetWeekDays() (weekDays []entity.WeekDay, err error) {
	return wdrs.repo.GetWeekDays()
}

func (wdrs BSWeekDayService) Get(bsID int, weekDayID int) (wds []entity.WeekDayBS, err error){
	return wdrs.repo.Get(bsID,weekDayID)
}