package bs

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/global/contracts"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/types"
	"bitbucket.org/lyndus/backend/infra/utils"
	"github.com/dgrijalva/jwt-go"
)

//LoginHandler login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user types.UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var service contracts.AuthServiceInterface
	service = services.NewBSAuthService()

	authentication, err := service.GetUserByCredential(user.CNPJ, user.Username)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if utils.ComparePasswords(authentication.Password, user.Password) {

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
	return
}

//RequestRecovery func recovery password
func RequestRecovery(w http.ResponseWriter, r *http.Request) {
	var user types.UserCredentials
	var service contracts.AuthServiceInterface

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		config.ResponsePerErr(w, err, config.DATAINVALID)
		return
	}

	service = services.NewBSAuthService()
	authentication, err := service.GetUserEmailBYUsername(user.CNPJ, user.Username)
	if err != nil {
		config.ResponsePerErr(w, err, config.USERNOTFOUND)
		return
	}
	if authentication.ID == 0 {
		config.ResponsePerErr(w, err, config.USERNOTFOUND)
		return
	}

	if !utils.ValidEmail(authentication.Email) {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	email := strings.Split(authentication.Email, "@")

	var respEmail string
	if len(email[0]) < 7 {
		respEmail = email[0][:2] + "***" + email[0][len(email[0])-2:] + "@" + email[1]
	}

	respEmail = email[0][:3] + "*****" + email[0][len(email[0])-3:] + "@" + email[1]

	cryptText := utils.Encrypt([]byte(strconv.Itoa(authentication.ID)), config.Config.SecretHash)

	response := types.UserRecover{
		Key:   cryptText,
		Email: respEmail,
	}

	config.JSONResponse(response, http.StatusOK, w)
	return
}

//ConfirmRecovery func confirm recovery password
func ConfirmRecovery(w http.ResponseWriter, r *http.Request) {
	var key types.Key

	err := json.NewDecoder(r.Body).Decode(&key)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	resp, err := utils.Decrypt(key.Key, config.Config.SecretHash)
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}
	id, err := strconv.Atoi(string(resp))
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	var authentication auth.Authentication
	var service contracts.AuthServiceInterface = services.NewBSAuthService()
	authentication, err = service.GetEmailBYID(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}
	if authentication.ID == 0 {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	code := utils.RandStringBytesMask(6)
	token := utils.HashPwd([]byte(code))

	err = auth.Service.ClearPasswordRecovery(authentication.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	_, err = auth.Service.CreatePasswordRecovery(authentication.ID, token)
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	tmpl, err := template.ParseFiles(config.Config.TemplatePath + "client_recovery_pass.html")
	if err != nil {

		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, code); err != nil {

		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return

	}

	result := tpl.String()
	email := utils.Email{
		To:      authentication.Email,
		Subject: "Recuperação de Senha",
		Body:    result,
	}
	emailConfig := utils.EmailConfig{
		Email:    config.Config.MailFrom,
		Password: config.Config.MailPassword,
		Host:     config.Config.MailSMTPHost,
		Port:     config.Config.MailSMTPPort,
	}

	err = utils.SendEmailSSL(emailConfig, email)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

//Recovery func recovery password
func Recovery(w http.ResponseWriter, r *http.Request) {
	var upr types.UserPasswordRecovery

	err := json.NewDecoder(r.Body).Decode(&upr)
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	resp, err := utils.Decrypt(upr.Key, config.Config.SecretHash)
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}
	id, err := strconv.Atoi(string(resp))
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	var authentication auth.Authentication
	var service contracts.AuthServiceInterface
	service = services.NewBSAuthService()

	authentication, err = service.GetTokenRecovery(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}
	if authentication.ID == 0 {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	if !utils.ComparePasswords(authentication.Token, upr.PasswordRecoveryCode) {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	newPass := utils.HashPwd([]byte(upr.NewPassword))

	err = auth.Service.UpdateAuthenticationsPassword(authentication.ID, newPass)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = auth.Service.ClearPasswordRecovery(authentication.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

//ChangePassword func change password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var changePassword types.ChangePassword

	id := r.Context().Value("id").(int)

	err := json.NewDecoder(r.Body).Decode(&changePassword)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var authentication auth.Authentication
	var service contracts.AuthServiceInterface
	service = services.NewBSAuthService()

	authentication, err = service.GetUserPasswordByID(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if authentication.ID == 0 {
		config.ResponsePerErr(w, err, config.USERNOTFOUND)
		return
	}

	if !utils.ComparePasswords(authentication.Password, changePassword.Password) {
		config.ResponsePerErr(w, err, config.INVALIDPASSWORD)
		return
	}

	newPass := utils.HashPwd([]byte(changePassword.NewPassword))

	err = auth.Service.UpdateAuthenticationsPassword(authentication.ID, newPass)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}


func ChangeUserPassword(w http.ResponseWriter, r *http.Request){

}