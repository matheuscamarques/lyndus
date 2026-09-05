package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BSDayOffService struct {
	repo contracts.BSDayOffRepositoryInterface
}

var DayOffService contracts.BSDayOffServiceInterface

func NewBSDayOffService() *BSDayOffService {
	return &BSDayOffService{
		repo: repository.NewBSDayOffRepository(),
	}
}

func (bdo BSDayOffService) AddBSDayOff(bsID int, day entity.DayOff) (id int, err error) {
	return bdo.repo.AddBSDayOff(bsID, day)
}
func (bdo BSDayOffService) GetBSDaysOff(bsID int) (days []entity.DayOff, err error) {
	return bdo.repo.GetBSDaysOff(bsID)
}
func (bdo BSDayOffService) DeleteBSDayOff(bsID, dayOffID int) error {
	return bdo.repo.DeleteBSDayOff(bsID, dayOffID)
}

func (bdo BSDayOffService) AddBSBMDayOff(bsID, bmID int, day entity.DayOff) (id int, err error) {
	return bdo.repo.AddBSBMDayOff(bsID, bmID, day)
}
func (bdo BSDayOffService) GetBSBMDaysOff(bsID, bmID int) (days []entity.DayOff, err error) {
	return bdo.repo.GetBSBMDaysOff(bsID, bmID)
}
func (bdo BSDayOffService) DeleteBSBMDayOff(bsID, bmID, dayOffID int) error {
	return bdo.repo.DeleteBSBMDayOff(bsID, bmID, dayOffID)
}
