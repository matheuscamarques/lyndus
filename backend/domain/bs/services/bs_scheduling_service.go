package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type BSScheduleService struct {
	repo contracts.BSScheduleRepositoryInterface
}

var Schedule contracts.BSScheduleServiceInterface

func NewBSSchedulingService() *BSScheduleService {
	return &BSScheduleService{
		repo: repository.NewBSScheduleRepository(),
	}
}

func (bss BSScheduleService) CreateScheduling(scheduling entity.Scheduling) (id int, err error) {
	return bss.repo.CreateScheduling(scheduling)
}

func (bss BSScheduleService) GetBSScheduling(bsId int, date time.Time) (scheduling []entity.Scheduling, err error) {
	return bss.repo.GetBSScheduling(bsId, date)
}

func (bss BSScheduleService) ChangeStatusSchedule(bsID, scheduleID, status int) (err error) {
	return bss.repo.ChangeStatusSchedule(bsID, scheduleID, status)
}

func (bss BSScheduleService) GetBSScheduledTime(bsId, bmID int, startTime, endTime types.DateTime) (ids []int, err error) {
	return bss.repo.GetBSScheduledTime(bsId, bmID, startTime, endTime)
}

func (bss BSScheduleService) GetBSSchedule(bsId, scheduleID int) (scheduling entity.Scheduling, err error) {
	return bss.repo.GetBSSchedule(bsId, scheduleID)
}

func (bss BSScheduleService) SetOrderIDSchedule(bsID, scheduleID, orderID int) (err error){
	return bss.repo.SetOrderIDSchedule(bsID, scheduleID, orderID)
}