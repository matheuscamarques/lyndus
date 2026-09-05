package contracts

import "bitbucket.org/lyndus/backend/domain/appuser/entity"

type AppUserStatementServiceInterface interface {
	GetStatements(appUserID int) (statements []entity.Statement, err error)
}
