package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type ClientAppUserRepository struct {
	conn *sqlx.DB
}

func NewClientAppUserRepository() ClientAppUserRepository {
	return ClientAppUserRepository{
		conn: postgres.DB,
	}
}

func (a ClientAppUserRepository) CreateAppUserBasic(user *entity.AppUser) error {
	command := `INSERT INTO app_user(person_id, cpf) VALUES($1, $2) RETURNING ID`
	stmt, err := a.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(user.PersonID, user.CPF).Scan(&user.ID)
	return err
}

func (a ClientAppUserRepository) GetAppUserIDByCPF(user *entity.AppUser) error {
	command := `SELECT aus.id FROM app_user aus WHERE aus.cpf=$1 `
	err := a.conn.Get(user, command, user.CPF)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
