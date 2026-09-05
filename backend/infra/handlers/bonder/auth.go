package bonder

import (
	"bitbucket.org/lyndus/backend/infra/auth"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/dgrijalva/jwt-go"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user types.UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	fakeAuth := auth.Authentication{}
	fakeAuth.ID = 99
	fakeAuth.Username = "lyndus"
	fakeAuth.Password = "lyndus@bonder2022"
	authentication := fakeAuth
	//var bs_service contracts.AuthServiceInterface = services.NewBonderAuthService()

	//authentication, err := bs_service.GetUserByCredential(bs_user.CNPJ, bs_user.Username)
	//if err != nil {
	//	config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	//	return
	//}

	if user.Username == authentication.Username && user.Password == authentication.Password {

		claims := jwt.StandardClaims{
			Id:        strconv.Itoa(authentication.ID),
			Subject:   authentication.Username,
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(config.Config.TokenExpires)).Unix(),
			IssuedAt:  time.Now().Unix(),
		}

		////log.Println(claims)
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		tokenString, err := token.SignedString(config.SignKey)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		response := types.UserResponse{
			AccessToken: tokenString,
			TokenType:   "RS256",
			Username:    authentication.Username,
			CompanyName: authentication.CompanyName,
		}

		config.JSONResponse(response, http.StatusOK, w)
		return
	}

	config.ResponsePerErr(w, err, config.INVALIDLOGIN)
}
