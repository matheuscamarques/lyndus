package repository

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type AppUserRepository struct {
	conn *sqlx.DB
}

func NewAppUserUserRepository() *AppUserRepository {
	return &AppUserRepository{
		conn: postgres.DB,
	}
}

func (aur AppUserRepository) UpdateAppUserPerson(appUser entity.AppUserCreate) error {
	command := `UPDATE bs_person SET  gender_id=$2,
       				   			   name=$3,
       				   			   birthdate=$4,
       				   			   phone=$5,
       				   			   cellphone=$6,
					   			   email=$7,
       				   			   state=$8,
       				   			   city=$9,
       				   			   district=$10,
       				   			   street=$11,
       				   			   number=$12,
       				   			   address_complement=$13,
       				   			   zipcode=$14
                        	WHERE id=$1`
	stmt, err := postgres.DB.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(appUser.PersonID,
		appUser.GenderID,
		appUser.Name,
		appUser.Birthdate,
		appUser.Phone,
		appUser.Cellphone,
		appUser.Email,
		appUser.State,
		appUser.City,
		appUser.District,
		appUser.Street,
		appUser.Number,
		appUser.AddressComplement,
		appUser.Zipcode)
	return err
}

func (aur AppUserRepository) UpdateAppUserEmail(appUser entity.AppUserCreate) error {
	command := `UPDATE app_user SET email =$3 
                        WHERE person_id=$1 AND id=$2`
	stmt, err := postgres.DB.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		appUser.PersonID,
		appUser.ID,
		appUser.Email)
	return err
}

func (aur AppUserRepository) UpdateAppUser(appUser entity.AppUser) error {
	command := `UPDATE app_user SET email =$3, 
                            		authentication_id =$4,
                            		accepted_term =$5
                        WHERE person_id=$1 AND id=$2`
	stmt, err := postgres.DB.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		appUser.PersonID,
		appUser.ID,
		appUser.Email,
		appUser.AuthenticationID,
		appUser.AcceptedTerm)
	return err
}

func (aur AppUserRepository) CreateAppUser(appUser entity.AppUser) (id int, err error) {
	command := `INSERT INTO app_user(person_id,
                     				 cpf,
                     				 email,
                     				 authentication_id,
                     				 accepted_term)
                                VALUES($1, $2, $3, $4, $5)
                                RETURNING ID`
	stmt, err := postgres.DB.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		appUser.PersonID,
		appUser.CPF,
		appUser.Email,
		appUser.AuthenticationID,
		appUser.AcceptedTerm).Scan(&id)
	return id, err
}

func (aur AppUserRepository) GetAppUserPersonIDByID(id int) (personID int, err error) {
	command := `SELECT aus.person_id
                FROM app_user aus 
                WHERE aus.id=$1 `
	err = aur.conn.Get(&personID, command, id)
	if err == sql.ErrNoRows {
		err = nil
	}
	return personID, err
}

func (aur AppUserRepository) GetAppUserByID(id int) (appUser entity.AppUserCreate, err error) {
	command := `SELECT aus.id,
                  	   per.name,
                  	   per.cellphone,
                  	   per.birthdate,
					   per.cpf
                FROM app_user aus 
                INNER JOIN bs_person per ON aus.person_id = per.id
                WHERE aus.id=$1 `
	err = aur.conn.Get(&appUser, command, id)
	if err == sql.ErrNoRows {
		err = nil
	}
	return appUser, err
}

func (aur AppUserRepository) GetAppUserFullByID(id int) (appUser entity.AppUserCreate, err error) {
	command := `SELECT aus.id,
       				   per.gender_id,
       				   per.name,
       				   aus.cpf,
       				   per.birthdate,
       				   per.phone,
       				   per.cellphone,
					   per.email,
       				   per.state,
       				   per.city,
       				   per.district,
       				   per.street,
       				   per.number,
       				   per.address_complement,
       				   per.zipcode
                FROM app_user aus 
                INNER JOIN bs_person per ON aus.person_id = per.id
                WHERE aus.id=$1 `
	err = aur.conn.Get(&appUser, command, id)
	if err == sql.ErrNoRows {
		err = nil
	}
	return appUser, err
}
