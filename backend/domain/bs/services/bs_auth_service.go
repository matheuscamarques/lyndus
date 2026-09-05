package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BSAuthService struct {
	repo contracts.BSAuthRepositoryInterface
}

var AuthService contracts.BSAuthServiceInterface

func NewBSAuthService() *BSAuthService {
	return &BSAuthService{
		repo: repository.NewBSAuthRepository(),
	}
}

//GetUserByCredential Get access by username
func (bas BSAuthService) GetUserByCredential(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	return bas.repo.GetUserByCredential(cnpj, username)
}

//GetUserEmailBYUsername Get access by username
func (bas BSAuthService) GetUserEmailBYUsername(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	return bas.repo.GetUserEmailBYUsername(cnpj, username)
}

//GetEmailBYID Get access email by id
func (bas BSAuthService) GetEmailBYID(id int) (a auth.Authentication, err error) {
	return bas.repo.GetEmailBYID(id)
}

//GetTokenRecovery token recovery by  id
func (bas BSAuthService) GetTokenRecovery(id int) (a auth.Authentication, err error) {
	return bas.repo.GetTokenRecovery(id)

}

//GetUserPasswordByID get password by id
func (bas BSAuthService) GetUserPasswordByID(id int) (a auth.Authentication, err error) {
	return bas.repo.GetUserPasswordByID(id)
}

//CheckPermission check if client bs_user have permission
func (bas BSAuthService) CheckPermission(authID, permissionID, levelID int) (id int, err error) {
	return bas.repo.CheckPermission(authID, permissionID, levelID)
}
