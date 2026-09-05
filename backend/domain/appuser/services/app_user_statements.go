package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser"
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
)

type AppUserStatementService struct {
	repo contracts.AppUserStatementRepositoryInterface
}

var StatementService contracts.AppUserStatementServiceInterface

func NewAppUserStatementService() *AppUserStatementService {
	return &AppUserStatementService{
		repo: repository.NewAppUserStatementRepository(),
	}
}

func (aub AppUserStatementService) GetStatements(appUserID int) (statements []entity.Statement, err error) {
	statements, err = aub.repo.GetStatements(appUserID)
	if err != nil {
		return
	}

	for k := range statements {
		statements[k].Name = appuser.Statements[statements[k].StatementsID]
		if statements[k].StatementsID == 2 {
			statements[k].Value = statements[k].Value * -1
		}
	}

	return
}
