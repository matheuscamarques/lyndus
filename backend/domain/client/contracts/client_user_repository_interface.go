package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/utils"
)

type ClientUserRepositoryInterface interface {
	Create(client entity.ClientUser) (int, error)
	GetAll(activePage, itemsPerPage, clientID, status int, search, orderBy string, sortDesc bool) (clientsUser []entity.ClientUser, totalItems, totalPages int, err error)
	Update(client entity.ClientUser) error
	Get(clientID int, clientUserID int) (clientUser entity.ClientUser, err error)
	//AddPermissions (cliente entity.ClientUser) (PermissionsID []utils.LevelPermission)
	AddPermissions(client entity.ClientUser)
	PermissionExist(cliUserID, levelId int) (bool, error)
	GenerateUsername(client entity.ClientUser) (string, error)
	UserPermissions(clientUserID int) (permissions []utils.LevelPermission, err error)
	GetMasterUserAuthID(id int) (int, error)
	GetAllV2(ctra criteria.Criteria) (response criteria.ClientClientUserResponse, err error)

	UpdateStatus(clientID, userID, statusID int) error
}
