package types

type UserCredentials struct {
	CNPJ     CNPJ   `json:"cnpj"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AppUserCredentials struct {
	Username CPF    `json:"username"`
	Password string `json:"password"`
}

type Key struct {
	Key []byte `json:"key"`
}

type UserRecover struct {
	Key   []byte `json:"key"`
	Email string `json:"email"`
}

type UserPasswordRecovery struct {
	Key                  []byte `json:"key"`
	NewPassword          string `json:"newPassword"`
	PasswordRecoveryCode string `json:"passwordRecoveryCode"`
}

type UserResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	Username    string `json:"username"`
	CompanyName string `json:"companyName"`
}

type ChangePassword struct {
	ID          int    `json:"id"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}
