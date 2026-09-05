package contracts

import "bitbucket.org/lyndus/backend/domain/appuser/entity"

type AppUserStatementRepositoryInterface interface {
	GetStatements(appUserID int) (statements []entity.Statement, err error)
}
