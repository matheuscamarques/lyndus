package entity

import (
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/types"
	"database/sql"
	"time"
)

type AppUser struct {
	ID               int       `json:"id" db:"id"`
	PersonID         int       `json:"personID" db:"person_id"`
	CPF              types.CPF `json:"cpf" db:"cpf"`
	Email            string    `json:"email" db:"email"`
	CreatedAT        time.Time `json:"createdAT" db:"created_at,omitempty"`
	AuthenticationID int       `json:"authenticationID" db:"authentication_id"`
	AcceptedTerm     bool      `json:"acceptedTerm" db:"accepted_term"`
}

//GetUserByCredential Get access by username
func GetUserByCredential(username types.CPF) (a auth.AuthenticationAppUser, err error) {
	command := `SELECT app.id,
       				   app.authentication_id,
					   per.name as username,
					   aut.password
                FROM app_user app
                INNER JOIN bs_person as per ON  app.person_id = per.id
                INNER JOIN authentication aut ON app.authentication_id = aut.id
                WHERE app.cpf = $1 AND aut.disabled IS FALSE`

	err = postgres.DB.Get(&a, command, username)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetAppUserUserEmailBYUsername Get access by username
func GetAppUserUserEmailBYUsername(username string) (a auth.AuthenticationAppUser, err error) {
	command := `SELECT app.id,
       				   app.authentication_id,
					   per.name as username,
					   per.email
                FROM app_user app
                INNER JOIN bs_person as per ON app.person_id=per.id 
                WHERE app.cpf = $1 AND authentication_id IS NOT NULL`

	err = postgres.DB.Get(&a, command, username)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetAppUserByID Get access by id
func GetAppUserByID(id int) (a auth.AuthenticationAppUser, err error) {
	command := `SELECT app.id,
                  aut.id as authentication_id,
                  per.name as username, 
                  aut.password  
                FROM app_user app
                INNER JOIN bs_person as per ON  app.person_id = per.id
                INNER JOIN authentication aut ON app.authentication_id = aut.id
                WHERE app.id=$1`

	err = postgres.DB.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetAppUserEmailBYID Get access email by id
func GetAppUserEmailBYID(id int) (a auth.AuthenticationAppUser, err error) {
	command := `SELECT app.id,
       				   app.authentication_id,
					   per.email
                FROM app_user app
                INNER JOIN bs_person as per ON app.person_id=per.id 
                WHERE app.id = $1`

	err = postgres.DB.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetTokenRecovery token recovery by app bs_user id
func GetTokenRecovery(id int) (a auth.AuthenticationAppUser, err error) {
	command := `SELECT app.id,
       				   app.authentication_id,
                       pre.token 
                FROM password_recovery pre 
                INNER JOIN app_user app ON pre.authentication_id = app.authentication_id
                WHERE app.id = $1 AND pre.created_at > now()::date - interval '1 hours'`

	err = postgres.DB.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetAppUserUserPasswordByID get app bs_user password by id
func GetAppUserUserPasswordByID(id int) (a auth.AuthenticationAppUser, err error) {
	command := `SELECT app.id,
       				   aut.id as authentication_id,
       				   aut.password
				FROM authentication aut
				INNER JOIN  app_user app on aut.id = app.authentication_id
				WHERE app.id = $1`

	err = postgres.DB.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}
