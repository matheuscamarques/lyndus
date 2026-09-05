package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
	_ "bitbucket.org/lyndus/backend/domain/client/repository"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/utils"
	"errors"
)

type ClientUserService struct {
	Repo contracts.ClientUserRepositoryInterface
}

var UserService contracts.ClientUserServiceInterface

func NewClientUserService() *ClientUserService {
	repo := repository.NewClientUserRepository()

	return &ClientUserService{
		Repo: &repo,
	}
}

func (ru ClientUserService) Update(user entity.ClientUser) error {
	//todo trocar password authentication table

	err := ru.Repo.Update(user)
	if err != nil {
		return err
	}

	if len(user.Permissions) > 0 {
		ru.Repo.AddPermissions(user)
	}

	return err
}

func (ru *ClientUserService) SetRepo(repo contracts.ClientUserRepositoryInterface) {
	ru.Repo = repo
}

func (ru ClientUserService) Get(clientID, clientUserID int) (entity.ClientUser, error) {
	return ru.Repo.Get(clientID, clientUserID)
}

func (ru ClientUserService) UpdateStatus(clientID, userID, statusID int) error {
	return ru.Repo.UpdateStatus(clientID, userID, statusID)
}

//func (ru ClientUserService) GetAll(clientID, status int) ([]entity.ClientUser, error) {
//	return ru.Repo.GetAll(clientID, status)
//}

func (ru ClientUserService) Register(user entity.ClientUser) (entity.ClientUser, error) {

	// Gera Username
	username, err := ru.Repo.GenerateUsername(user)
	if err != nil {
		return entity.ClientUser{}, err
	}
	user.Username = username
	user.StatusID = config.USERSTATUSACTIVE
	// Registra Usuario
	id, err := ru.Repo.Create(user)
	if err != nil {
		return entity.ClientUser{}, err
	}

	if user.Master {
		user.Permissions = utils.PermissionClientMaster
	}
	// Adiciona as permissões necessárias
	if id != 0 {
		user.ID = id
		if len(user.Permissions) > 0 {
			ru.Repo.AddPermissions(user)
		}
		return user, nil
	}

	return entity.ClientUser{}, errors.New("falha ao criar o usuário")
}

func (ru ClientUserService) UserPermissions(clientUserID int) (permissions []utils.LevelPermission, err error) {
	return ru.Repo.UserPermissions(clientUserID)
}

func (ru ClientUserService) GetAllV2(activePage, itemsPerPage, clientID int, search, orderBy string, sortDesc, active bool) (clientUsers []entity.ClientUser, totalItems, totalPages int, err error) {
	if activePage == 0 {
		activePage = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 10
	} else if itemsPerPage > 100 {
		itemsPerPage = 100
	}
	status := 2
	if active {
		status = 1
	}

	clientUsers, totalItems, totalPages, err = ru.Repo.GetAll(activePage, itemsPerPage, clientID, status, search, orderBy, sortDesc)
	if clientUsers == nil {
		clientUsers = make([]entity.ClientUser, 0)
	}

	return clientUsers, totalItems, totalPages, err
}
