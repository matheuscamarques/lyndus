package services

import (
	"bitbucket.org/lyndus/backend/global/aggregate"
	"time"

	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
)

type AppUserBMService struct {
	repo contracts.AppUserBMRepositoryInterface
}

var BM contracts.AppUserBMInterface

func NewAppUserBMService() *AppUserBMService {
	return &AppUserBMService{
		repo: repository.NewAppUserBMRepository(),
	}
}

func (ass AppUserBMService) GetBMsBasicByServiceID(bsID, serviceID int) (bms []entity.BM, err error) {
	return ass.repo.GetBMsBasicByServiceID(bsID, serviceID)
}

func (ass AppUserBMService) GetBMsByServiceID(bsID, serviceID int) (bms []entity.BM, err error) {
	return ass.repo.GetBMsByServiceID(bsID, serviceID)
}

func (ass AppUserBMService) GetBMWeekDays(bmID int) (weekDays []entity.WeekDayBS, err error) {
	return ass.repo.GetBMWeekDays(bmID)
}

func (ass AppUserBMService) GetBMByID(bsID, bmID int) (bm entity.BM, err error) {
	return ass.repo.GetBMByID(bsID, bmID)
}

func (ass AppUserBMService) GetBMByIDAndDuration(bsID, bmID, serviceID int) (bm entity.BM, err error) {
	return ass.repo.GetBMByIDAndDuration(bsID, bmID, serviceID)
}
func (ass AppUserBMService) GetBMDaysOff(bsID, bmID int, start, end time.Time) (days []entity.DayOff, err error) {
	return ass.repo.GetBMDaysOff(bsID, bmID, start, end)
}
func (ass AppUserBMService) GetBMsFullByWeekDay(bsID, serviceID, weekDayID int) (bm []entity.BM, err error) {
	return ass.repo.GetBMsFullByWeekDay(bsID, serviceID, weekDayID)
}
func (ass AppUserBMService) GetBMFullByWeekDay(bsID, serviceID, bmID, weekDayID int) (bm entity.BM, err error) {
	return ass.repo.GetBMFullByWeekDay(bsID, serviceID, bmID, weekDayID)
}

func (ass AppUserBMService) GetBMWeekDaysWithName(bmID int) (weekDays []aggregate.WeekDay, err error) {
	return ass.repo.GetBMWeekDaysWithName(bmID)
}
