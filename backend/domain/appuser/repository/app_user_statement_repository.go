package repository

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type AppUserStatementRepository struct {
	conn *sqlx.DB
}

func NewAppUserStatementRepository() AppUserStatementRepository {
	return AppUserStatementRepository{
		conn: postgres.DB,
	}
}

func (ausr AppUserStatementRepository) GetStatements(appUserID int) (statements []entity.Statement, err error) {
	command := `SELECT value, date_time, statements_id FROM app_user_statement WHERE app_user_id=$1`
	err = ausr.conn.Select(&statements, command, appUserID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}
