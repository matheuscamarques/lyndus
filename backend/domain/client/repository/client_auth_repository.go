package repository

import (
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/types"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type ClientAuthRepository struct {
	conn *sqlx.DB
}

func NewClientAuthRepository() ClientAuthRepository {
	return ClientAuthRepository{
		conn: postgres.DB,
	}
}

//GetUserByCredential Get access by username
func (cur ClientAuthRepository) GetUserByCredential(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	command := `SELECT aut.id as id,
					  com.company_name,
       				  clu.name as username,
					  aut.password
					  FROM client cli
					  INNER JOIN client_user clu ON cli.id = clu.client_id
					  INNER JOIN company as com ON cli.company_id=com.id 
					  INNER JOIN authentication aut ON clu.authentication_id = aut.id
					  WHERE cli.active_lyndus = true
					    AND deleted = false
					    AND cli.cnpj = $1
					    AND clu.username = $2
					    AND clu.status = 1
					    AND aut.disabled IS FALSE`

	err = cur.conn.Get(&a, command, cnpj, username)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetUserEmailBYUsername Get access by username
func (cur ClientAuthRepository) GetUserEmailBYUsername(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	command := `SELECT aut.id,
                  clu.name as username,                
       			  clu.email
            FROM client_user clu 
            INNER JOIN client cli ON cli.id = clu.client_id
			INNER JOIN authentication aut ON clu.authentication_id = aut.id
            WHERE cli.cnpj = $1 AND clu.username=$2 AND clu.status = 1 AND aut.disabled IS FALSE`

	err = cur.conn.Get(&a, command, cnpj, username)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetEmailBYID Get access email by id
func (cur ClientAuthRepository) GetEmailBYID(id int) (a auth.Authentication, err error) {
	command := `SELECT clu.authentication_id as id,
                  clu.email
                FROM client_user clu 
                WHERE clu.authentication_id = $1`

	err = cur.conn.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetTokenRecovery token recovery by id
func (cur ClientAuthRepository) GetTokenRecovery(id int) (a auth.Authentication, err error) {
	command := `SELECT  clu.authentication_id as id,
						pre.token
					FROM password_recovery pre
					INNER JOIN client_user clu ON pre.authentication_id = clu.authentication_id
					WHERE clu.authentication_id = $1  AND pre.created_at > now()::date - interval '1 hours'`

	err = cur.conn.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetUserPasswordByID get client password by id
func (cur ClientAuthRepository) GetUserPasswordByID(id int) (a auth.Authentication, err error) {
	command := `SELECT clu.authentication_id as id,
       				   aut.password
				FROM authentication aut
				INNER JOIN client_user clu  ON clu.authentication_id = aut.id
				WHERE clu.authentication_id = $1`

	err = cur.conn.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, err
	}
	return a, err
}

//CheckPermission check if client bs_user have permission
func (cur ClientAuthRepository) CheckPermission(authID, permissionID, levelID int) (clientID int, err error) {
	//TODO talvez não se encaixe como repository
	command := `SELECT clu.client_id
				FROM client_user clu
				INNER JOIN client_user_permission cup ON cup.client_user_id = clu.id
				WHERE clu.authentication_id = $1
				  AND cup.client_permission_id = $2
				  AND access_level >= $3`

	err = cur.conn.Get(&clientID, command, authID, permissionID, levelID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return clientID, err
}
