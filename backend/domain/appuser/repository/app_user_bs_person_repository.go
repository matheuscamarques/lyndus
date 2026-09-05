package repository

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/types"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
)

type AppUserBSPersonRepository struct {
	conn *sqlx.DB
}

func NewAppUserBSPersonRepository() *AppUserBSPersonRepository {
	return &AppUserBSPersonRepository{
		conn: postgres.DB,
	}
}

func (aub AppUserBSPersonRepository) GetBSPersonByCPF(bsID int, cpf types.CPF) (bm entity.BSPerson, err error) {
	command := `SELECT bsp.id,
       				   bsp.app_user_id,
       				   bsp.cpf
					FROM  bs_person bsp
					WHERE bsp.bs_id = $1 
					AND   bsp.cpf = $2 `
	err = aub.conn.Get(&bm, command, bsID, cpf)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bm, err
}

func (aub AppUserBSPersonRepository) GetBSPersonByAppUserID(bsID, appUserID int) (bm entity.BSPerson, err error) {
	command := `SELECT bsp.id,
       				   bsp.app_user_id,       				   
       				   bsp.cpf
					FROM  bs_person bsp
					WHERE bsp.bs_id = $1 
					AND   bsp.app_user_id = $2 `
	err = aub.conn.Get(&bm, command, bsID, appUserID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bm, err
}

func (aub AppUserBSPersonRepository) CreatePerson(person entity.BSPerson) (id int, err error) {
	log.Println(person)
	command := `INSERT INTO bs_person(bs_id,
                      				  app_user_id,
                      				  name,
                      				  cpf,
                      				  phone,
                      				  birthdate) 
            		VALUES($1, $2, $3, $4, $5, $6) RETURNING ID`
	stmt, err := aub.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(person.BsID,
		person.AppUserID,
		person.Name,
		person.CPF,
		person.Phone,
		person.Birthdate).Scan(&id)
	return id, err
}

func (aub AppUserBSPersonRepository) UpdatePersonAppUserID(bsID, personID, appUserID int) error {
	command := `UPDATE bs_person SET app_user_id=$3 
								 WHERE bs_id=$1
									AND id=$2`
	stmt, err := aub.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, personID, appUserID)
	return err
}
