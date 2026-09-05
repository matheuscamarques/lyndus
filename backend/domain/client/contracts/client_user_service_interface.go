package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/utils"
)

type ClientUserServiceInterface interface {
	Register(user entity.ClientUser) (entity.ClientUser, error)
	Update(user entity.ClientUser) error
	Get(clientID int, clientUserID int) (entity.ClientUser, error)
	//GetAll(userid, status int) ([]entity.ClientUser, error)
	UserPermissions(clientUserID int) (permissions []utils.LevelPermission, err error)
	//GetAllV2(id int, status int, page int, itemsPerPage int, search string) (response criteria.ClientClientUserResponse, err error)
	GetAllV2(activePage, itemsPerPage, clientID int, search, orderBy string, sortDesc, active bool) ([]entity.ClientUser, int, int, error)

	UpdateStatus(clientID, userID, statusID int) error
}
