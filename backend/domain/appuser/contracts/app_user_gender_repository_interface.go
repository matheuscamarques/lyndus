package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
)

type AppUserGenderRepositoryInterface interface {
	GetGenders() (genders []entity.Gender, err error)
	GetCategories() (categories []entity.Category, err error)
	GetGenderID(genderID int) (id int, err error)
	GetGenderByID(genderID int) (gender entity.Gender, err error)
}
