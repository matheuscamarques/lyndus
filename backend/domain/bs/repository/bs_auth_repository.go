package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/types"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type BSAuthRepository struct {
	contracts.BSAuthRepositoryInterface
	conn *sqlx.DB
}

func NewBSAuthRepository() BSAuthRepository {
	return BSAuthRepository{
		conn: postgres.DB,
	}
}

//GetUserByCredential Get access by username
func (bar BSAuthRepository) GetUserByCredential(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	command := `SELECT  aut.id,
						com.company_name, 
						bsu.name as username,
						aut.password
					FROM authentication aut 
					INNER JOIN bs_user bsu ON aut.id = bsu.authentication_id
					INNER JOIN bs ON bsu.bs_id = bs.id
					INNER JOIN company com ON bs.company_id=com.id 
					WHERE com.cnpj = $1 AND bsu.username = $2 AND bsu.status = 1 AND aut.disabled IS FALSE`

	err = bar.conn.Get(&a, command, cnpj, username)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetUserEmailBYUsername Get access by username
func (bar BSAuthRepository) GetUserEmailBYUsername(cnpj types.CNPJ, username string) (a auth.Authentication, err error) {
	command := `SELECT aut.id,
       				   bsu.name as username,
       				   bsu.email
            FROM bs_user bsu
            INNER JOIN authentication aut on bsu.authentication_id = aut.id
			INNER JOIN bs ON bsu.bs_id = bs.id
            WHERE bs.cnpj = $1 AND bsu.username = $2 AND bsu.status = 1 AND aut.disabled IS FALSE`

	err = bar.conn.Get(&a, command, cnpj, username)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetEmailBYID Get access email by id
func (bar BSAuthRepository) GetEmailBYID(id int) (a auth.Authentication, err error) {
	command := `SELECT bsu.authentication_id as id,
                  bsu.email
                FROM bs_user bsu
                WHERE bsu.authentication_id = $1`

	err = bar.conn.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetTokenRecovery token recovery by  id
func (bar BSAuthRepository) GetTokenRecovery(id int) (a auth.Authentication, err error) {
	command := `SELECT bsu.authentication_id as id,
						pre.token
					FROM password_recovery pre
					INNER JOIN bs_user bsu  ON pre.authentication_id = bsu.authentication_id
					WHERE bsu.authentication_id = $1  AND pre.created_at > now()::date - interval '1 hours'`

	err = bar.conn.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, nil
	}
	return a, err
}

//GetUserPasswordByID get password by id
func (bar BSAuthRepository) GetUserPasswordByID(id int) (a auth.Authentication, err error) {
	command := `SELECT aut.id,
       				   aut.password
				FROM authentication aut
				INNER JOIN bs_user bsu ON bsu.authentication_id = aut.id
				WHERE aut.id = $1`

	err = bar.conn.Get(&a, command, id)
	if err == sql.ErrNoRows {
		return a, err
	}
	return a, err
}

//CheckPermission check if bs bs_user have permission
func (bar BSAuthRepository) CheckPermission(authID, permissionID, levelID int) (bsID int, err error) {
	command := `SELECT bsu.bs_id
				FROM bs_user bsu
				INNER JOIN bs_user_permission bup ON bsu.id = bup.bs_user_id
				WHERE bsu.authentication_id = $1
				  AND bup.bs_permission_id=$2
				  AND bup.access_level >= $3`

	err = bar.conn.Get(&bsID, command, authID, permissionID, levelID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return bsID, err
}
