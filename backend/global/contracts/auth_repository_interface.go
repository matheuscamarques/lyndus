package contracts

import (
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/types"
)

type AuthRepositoryInterface interface {
	GetUserByCredential(cnpj types.CNPJ, username string) (a auth.Authentication, err error)
	GetUserEmailBYUsername(cnpj types.CNPJ, username string) (a auth.Authentication, err error)
	GetEmailBYID(id int) (a auth.Authentication, err error)
	GetTokenRecovery(id int) (a auth.Authentication, err error)
	GetUserPasswordByID(id int) (a auth.Authentication, err error)
	CheckPermission(authID, permissionID, levelID int) (clientID int, err error)
}
