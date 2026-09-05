package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/jmoiron/sqlx"
)

type BSAppUserRepository struct {
	conn *sqlx.DB
}

func NewBSAppUserRepository() BSAppUserRepository {
	return BSAppUserRepository{
		conn: postgres.DB,
	}
}

func (a BSAppUserRepository) GetAppUserByCPF(cpf types.CPF) (appUser entity.AppUser, err error) {
	command := `SELECT aus.id,
                  	   per.name,
                  	   per.cellphone,
                  	   per.birthdate
                FROM app_user aus 
                INNER JOIN bs_person per ON aus.person_id = per.id
                WHERE aus.cpf=$1 `
	err = a.conn.Get(appUser, command, cpf)
	if err == sql.ErrNoRows {
		err = nil
	}
	return appUser, err
}

func (a BSAppUserRepository) GetAppUserIDByCPF(cpf *types.CPF) (id int, err error) {
	command := `SELECT aus.id FROM app_user aus WHERE aus.cpf=$1`
	err = a.conn.Get(&id, command, cpf)
	if err == sql.ErrNoRows {
		err = nil
	}
	return id, err
}
