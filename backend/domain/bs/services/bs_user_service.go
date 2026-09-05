package services

import (
	"bitbucket.org/lyndus/backend/infra/utils"
	"errors"

	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/infra/config"
)

type BSUserService struct {
	Repo contracts.BSUserRepositoryInterface
}

var UserService contracts.BSUserServiceInterface

func NewBSUserService() *BSUserService {
	return &BSUserService{
		Repo: repository.NewBSUserRepository(),
	}
}

func (ru *BSUserService)SetRepo(repo contracts.BSUserRepositoryInterface){
	ru.Repo = repo
}
func (ru BSUserService) Update(bsUser entity.BSUser) error {

	//todo trocar password authentication table
	err := ru.Repo.Update(bsUser)
	if err != nil {
		return err
	}

	//update permissions
	if bsUser.ID != 0 {
		if len(bsUser.Permissions) > 0 {
			ru.Repo.AddPermissions(bsUser)
		}
	}
	return err
}

func (ru BSUserService) Get(bsID, BSUserID int) (entity.BSUser, error) {
	return ru.Repo.Get(bsID, BSUserID)
}

func (ru BSUserService) GetAll(bsID, status int) ([]entity.BSUser, error) {
	return ru.Repo.GetAll(bsID, status)
}

func (ru BSUserService) Register(bsUser entity.BSUser) (entity.BSUser, error) {

	// Gera Username
	username, err := ru.Repo.GenerateUsername(bsUser)
	if err != nil {
		return entity.BSUser{}, err
	}
	bsUser.Username = username
	bsUser.Status = config.USERSTATUSACTIVE
	// Registra Usuário
	id, err := ru.Repo.Create(bsUser)

	if err != nil {
		return entity.BSUser{}, err
	}

	if bsUser.Master {
		bsUser.Permissions = utils.PermissionBsMaster
	}
	// Adiciona as permissões necessárias
	if id != 0 {
		bsUser.ID = id
		if len(bsUser.Permissions) > 0 {
			ru.Repo.AddPermissions(bsUser)
		}
		return bsUser, err
	}

	return entity.BSUser{}, errors.New("falha ao criar o usuário")
}
