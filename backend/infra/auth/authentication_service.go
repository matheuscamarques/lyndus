package auth

import (
	"time"

	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/utils"
	"github.com/jmoiron/sqlx"
)

type AuthenticationService struct {
	conn *sqlx.DB
}

var Service *AuthenticationService

// TODO MAKE INTERFACE FOR THIS AUTHENTICATOR
func NewAuthenticationService() *AuthenticationService {
	return &AuthenticationService{
		conn: postgres.DB,
	}
}

//CreateAuthentication save password recovery
func (au AuthenticationService) CreateAuthentication(password string) (n int, err error) {

	// create password hash
	hashPassword := utils.HashPwd([]byte(password))

	sqlCommand := `INSERT INTO authentication(password) VALUES($1) RETURNING ID`
	stmt, err := au.conn.Prepare(sqlCommand)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(hashPassword).Scan(&n)
	return
}

//CreatePasswordRecovery save password recovery
func (au AuthenticationService) CreatePasswordRecovery(authID int, token string) (n int, err error) {
	sqlCommand := `INSERT INTO password_recovery(authentication_id, token) VALUES($1, $2) RETURNING ID`
	stmt, err := au.conn.Prepare(sqlCommand)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(authID, token).Scan(&n)
	return
}

//UpdateAuthenticationsPassword update client password
func (au AuthenticationService) UpdateAuthenticationsPassword(authID int, password string) error {
	sqlCommand := `UPDATE authentication SET password=$2 WHERE id=$1`
	stmt, err := au.conn.Prepare(sqlCommand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(authID, password)
	return err
}

//ClearPasswordRecovery clear password recovery
func (au AuthenticationService) ClearPasswordRecovery(authID int) error {
	sqlCommand := `DELETE FROM password_recovery WHERE authentication_id = $1`
	stmt, err := au.conn.Prepare(sqlCommand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(authID)
	return err
}

func (au *AuthenticationService) CreateAuthToken(authenticationID int, expiration time.Time) (n int, err error) {
	command := `INSERT INTO auth_token(authentication_id, expiration) VALUES($1, $2) RETURNING ID`
	stmt, err := au.conn.Prepare(command)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(authenticationID, expiration).Scan(&n)
	return
}
