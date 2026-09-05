package services

import (
	"bitbucket.org/lyndus/backend/domain/bonder/contracts"
	"bitbucket.org/lyndus/backend/domain/bonder/repository"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BonderAuthService struct {
	repo contracts.BonderAuthRepositoryInterface
}

var AuthService contracts.BonderAuthServiceInterface

func NewBonderAuthService() *BonderAuthService {
	return &BonderAuthService{
		repo: repository.NewBonderAuthRepository(),
	}
}

//GetUserByCredential Get access by username
func (bas BonderAuthService) GetUserByCredential(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	return bas.repo.GetUserByCredential(cnpj, username)
}

//GetUserEmailBYUsername Get access by username
func (bas BonderAuthService) GetUserEmailBYUsername(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	return bas.repo.GetUserEmailBYUsername(cnpj, username)
}

//GetEmailBYID Get access email by id
func (bas BonderAuthService) GetEmailBYID(id int) (a auth.Authentication, err error) {
	return bas.repo.GetEmailBYID(id)
}

//GetTokenRecovery token recovery by  id
func (bas BonderAuthService) GetTokenRecovery(id int) (a auth.Authentication, err error) {
	return bas.repo.GetTokenRecovery(id)

}

//GetUserPasswordByID get password by id
func (bas BonderAuthService) GetUserPasswordByID(id int) (a auth.Authentication, err error) {
	return bas.repo.GetUserPasswordByID(id)
}

//CheckPermission check if client bs_user have permission
func (bas BonderAuthService) CheckPermission(authID, permissionID, levelID int) (id int, err error) {
	return bas.repo.CheckPermission(authID, permissionID, levelID)
}
