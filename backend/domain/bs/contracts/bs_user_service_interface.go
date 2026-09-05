package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSUserServiceInterface interface {
	Register(bsUser entity.BSUser) (entity.BSUser,  error)
	Update(bsUser entity.BSUser) error
	Get(bsId int, BSUserID int) (entity.BSUser, error)
	GetAll(bsID, status int) ([]entity.BSUser, error)
}
