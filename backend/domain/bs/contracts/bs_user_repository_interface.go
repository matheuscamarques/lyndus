package contracts

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
)

type BSUserRepositoryInterface interface {
	Create(bsUser entity.BSUser) (int,error)
	GetAll(bsID, status int) ([]entity.BSUser,error)
	Update(bsUser entity.BSUser) error
	Get(bsID int, BSUserID int) (BSUser entity.BSUser, err error)
	// AddPermissions (cliente entity.BSUser) (PermissionsID []utils.LevelPermission)
	AddPermissions(bsUser entity.BSUser) ()
	PermissionExist(bsUserID, levelId int)(bool , error)
	GenerateUsername(bsUser entity.BSUser) (string,error)
	GetMasterUserAuthID(id int) (int,error)
}