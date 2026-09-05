package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/repository"
	"bitbucket.org/lyndus/backend/infra/auth"

	"bitbucket.org/lyndus/backend/infra/types"
)

type ClientAuthService struct {
	repo contracts.ClientAuthRepositoryInterface
}

var AuthService contracts.ClientAuthServiceInterface

func NewClientAuthService() *ClientAuthService {
	repo := repository.NewClientAuthRepository()
	return &ClientAuthService{
		repo: &repo,
	}
}

//GetUserByCredential Get access by username
func (cus ClientAuthService) GetUserByCredential(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	return cus.repo.GetUserByCredential(cnpj, username)
}

//GetUserEmailBYUsername Get access by username
func (cus ClientAuthService) GetUserEmailBYUsername(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	return cus.repo.GetUserEmailBYUsername(cnpj, username)
}

//GetEmailBYID Get access email by id
func (cus ClientAuthService) GetEmailBYID(id int) (a auth.Authentication, err error) {
	return cus.repo.GetEmailBYID(id)
}

//GetTokenRecovery token recovery by id
func (cus ClientAuthService) GetTokenRecovery(id int) (a auth.Authentication, err error) {
	return cus.repo.GetTokenRecovery(id)
}

//GetUserPasswordByID get client password by id
func (cus ClientAuthService) GetUserPasswordByID(id int) (a auth.Authentication, err error) {
	return cus.repo.GetUserPasswordByID(id)
}

//CheckPermission check if client bs_user have permission
func (cus ClientAuthService) CheckPermission(authID, permissionID, levelID int) (id int, err error) {
	return cus.repo.CheckPermission(authID, permissionID, levelID)
}
