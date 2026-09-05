package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
	"bitbucket.org/lyndus/backend/global/basic"
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type AppUserScheduleService struct {
	repo contracts.AppUserScheduleRepositoryInterface
}

var Schedule contracts.AppUserScheduleServiceInterface

func NewAppUserSchedulingService() *AppUserScheduleService {
	return &AppUserScheduleService{
		repo: repository.NewAppUserScheduleRepository(),
	}
}

func (bss AppUserScheduleService) CreateScheduling(scheduling entity.Scheduling) (id int, err error) {
	return bss.repo.CreateScheduling(scheduling)
}

func (bss AppUserScheduleService) GetBSScheduling(bsId int, date time.Time) (scheduling []entity.Scheduling, err error) {
	return bss.repo.GetBSScheduling(bsId, date)
}

func (bss AppUserScheduleService) ChangeStatusSchedule(bsID, scheduleID, status int) (err error) {
	return bss.repo.ChangeStatusSchedule(bsID, scheduleID, status)
}

func (bss AppUserScheduleService) GetBSScheduledTime(bsId, bmID int, startTime, endTime types.DateTime) (ids []int, err error) {
	return bss.repo.GetBSScheduledTime(bsId, bmID, startTime, endTime)
}

func (bss AppUserScheduleService) GetBSBMScheduling(bsId, bmID int, start, end time.Time) (scheduling []entity.SchedulingBAsic, err error) {
	return bss.repo.GetBSBMScheduling(bsId, bmID, start, end)
}

func (bss AppUserScheduleService) GetSchedulingsByAppUserID(appUserID, activePage, itemsPerPage int, date types.Date) (scheduling []basic.AppScheduling, totalPages int, err error) {
	return bss.repo.GetSchedulingsByAppUserID(appUserID, activePage, itemsPerPage, date)
}

func (bss AppUserScheduleService) GetSchedulingAppUser(appUserID, scheduleID int) (scheduling entity.Scheduling, err error) {
	return bss.repo.GetSchedulingAppUser(appUserID, scheduleID)
}
