package contracts

import "bitbucket.org/lyndus/backend/domain/appuser/entity"

type AppUserGenderServiceInterface interface {
	GetGenders() (genders []entity.Gender, err error)
	GetCategories() (categories []entity.Category, err error)
	GetGenderID(id int) (int, error)
	GetGenderByID(genderID int) (gender entity.Gender, err error)
}
