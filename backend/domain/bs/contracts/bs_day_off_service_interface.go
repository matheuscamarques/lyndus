package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSDayOffServiceInterface interface {
	AddBSDayOff(bsID int, day entity.DayOff) (id int, err error)
	GetBSDaysOff(bsID int) (days []entity.DayOff, err error)
	DeleteBSDayOff(bsID, dayOffID int) error

	AddBSBMDayOff(bsID, bmID int, day entity.DayOff) (id int, err error)
	GetBSBMDaysOff(bsID, bmID int) (days []entity.DayOff, err error)
	DeleteBSBMDayOff(bsID, bmID, dayOffID int) error
}
