package auth

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"bitbucket.org/lyndus/backend/infra/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type Authentication struct {
	ID          int    `json:"id" db:"id""`
	CompanyName string `json:"companyName,omitempty" db:"company_name"`
	Username    string `json:"username" db:"username"`
	//FullName         string `json:"full_name" db:"full_name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`

	Token string `json:"token,omitempty" db:"token"`
	//RegisterDate time.Time  `json:"register_date,omitempty" db:"register_date"`
	//Deleted      bool       `json:"deleted,omitempty"`
	//LastUpdate   *time.Time `json:"last_update,omitempty" db:"last_update"`
}

type AuthenticationAppUser struct {
	Authentication
	AuthenticationID int `json:"authentication_id,omitempty" db:"authentication_id"`
}

//ValidateTokenMiddleware Validate Token request
func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return config.VerifyKey, nil
			})
		if err != nil {
			config.ResponsePerErr(
				w,
				errors.New("invalid token: "+err.Error()),
				config.UNAUTHORIZED)
			return
		}
		if !token.Valid {
			config.ResponsePerErr(
				w,
				errors.New("invalid token: "),
				config.UNAUTHORIZED)
			return
		}
		//todo jwt jti black list
		id, ok := token.Claims.(jwt.MapClaims)["jti"].(string)
		userID, _ := strconv.Atoi(id)

		if !ok || userID < 1 {
			//log.Println("Invalid jti")
			config.ResponsePerErr(w, errors.New("\"Invalid jti.\""), config.UNAUTHORIZED)
			return
		}
		ctx := context.WithValue(r.Context(), "id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func InitKeysPem() error {
	// TODO VERTIFICAR LOCAL DA CHAVE PRIVADA
	signBytes, err := ioutil.ReadFile(config.Config.PrivKeyPath)
	if err != nil {
		return err
	}

	config.SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}
	verifyBytes, err := ioutil.ReadFile(config.Config.PubKeyPath)
	if err != nil {
		return err
	}

	config.VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}
	return nil
}
