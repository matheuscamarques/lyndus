package contracts

import "time"

type AuthenticationServiceInterface interface {
	CreateAuthentication(password string) (n int, err error)
	CreatePasswordRecovery(authID int, token string) (n int, err error)
	UpdateAuthenticationsPassword(authID int, password string) error
	ClearPasswordRecovery(authID int) error
	CreateAuthToken(authenticationID int, expiration time.Time) (n int, err error)
}
