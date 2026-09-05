package appuser

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/types"
	"bitbucket.org/lyndus/backend/infra/utils"
	"github.com/dgrijalva/jwt-go"
)

//LoginHandler login func
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user types.AppUserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var access auth.AuthenticationAppUser

	access, err = entity.GetUserByCredential(user.Username)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if utils.ComparePasswords(access.Password, user.Password) {
		claims := jwt.StandardClaims{
			Id:        strconv.Itoa(access.ID),
			Subject:   access.Username,
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(config.Config.TokenExpires)).Unix(),
			IssuedAt:  time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		tokenString, err := token.SignedString(config.SignKey)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		response := types.UserResponse{
			AccessToken: tokenString,
			TokenType:   "RS256",
			Username:    access.Username,
		}

		config.JSONResponse(response, http.StatusOK, w)
		return
	}
	config.ResponsePerErr(w, err, config.INVALIDLOGIN)
}

//RequestRecovery func recovery password
func RequestRecovery(w http.ResponseWriter, r *http.Request) {
	var user types.UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		config.ResponsePerErr(w, err, config.DATAINVALID)
		return
	}

	var access auth.AuthenticationAppUser
	access, err = entity.GetAppUserUserEmailBYUsername(user.Username)
	if err != nil {
		config.ResponsePerErr(w, err, config.USERNOTFOUND)
		return
	}
	log.Println(access, user.Username, err)
	if access.ID == 0 {
		config.ResponsePerErr(w, err, config.USERNOTFOUND)
		return
	}

	if !utils.ValidEmail(access.Email) {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	email := strings.Split(access.Email, "@")

	var respEmail string
	if len(email[0]) < 7 {
		respEmail = email[0][:2] + "***" + email[0][len(email[0])-2:] + "@" + email[1]
	}

	respEmail = email[0][:3] + "*****" + email[0][len(email[0])-3:] + "@" + email[1]

	cryptText := utils.Encrypt([]byte(strconv.Itoa(access.ID)), config.Config.SecretHash)

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

	var access auth.AuthenticationAppUser
	access, err = entity.GetAppUserEmailBYID(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}
	if access.ID == 0 {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	code := utils.RandStringBytesMask(6)
	token := utils.HashPwd([]byte(code))

	err = auth.Service.ClearPasswordRecovery(access.AuthenticationID)
	if err != nil {
		config.ResponsePerErr(w, err, config.EXPIRED)
		return
	}

	_, err = auth.Service.CreatePasswordRecovery(access.AuthenticationID, token)
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
		To:      access.Email,
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
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if string(upr.Key) == "" || len(upr.PasswordRecoveryCode) < 4 || len(upr.NewPassword) < 3 {
		//log.Println("Recovery vazio:", upr)
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	resp, err := utils.Decrypt(upr.Key, config.Config.SecretHash)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	id, err := strconv.Atoi(string(resp))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var access auth.AuthenticationAppUser
	access, err = entity.GetTokenRecovery(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if access.AuthenticationID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if !utils.ComparePasswords(access.Token, upr.PasswordRecoveryCode) {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	newPass := utils.HashPwd([]byte(upr.NewPassword))

	err = auth.Service.UpdateAuthenticationsPassword(access.AuthenticationID, newPass)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = auth.Service.ClearPasswordRecovery(access.AuthenticationID)
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

	id, ok := r.Context().Value("id").(int)
	if !ok {
		config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&changePassword)
	if err != nil {
		config.ResponsePerErr(w, nil, config.EXPIRED)
		return
	}

	var access auth.AuthenticationAppUser

	access, err = entity.GetAppUserUserPasswordByID(id)
	if err != nil {
		config.ResponsePerErr(w, nil, config.EXPIRED)
		return
	}

	if access.AuthenticationID == 0 {
		config.ResponsePerErr(w, nil, config.USERNOTFOUND)
		return
	}

	if !utils.ComparePasswords(access.Password, changePassword.Password) {
		config.ResponsePerErr(w, nil, config.INVALIDPASSWORD)
		return
	}

	newPass := utils.HashPwd([]byte(changePassword.NewPassword))

	err = auth.Service.UpdateAuthenticationsPassword(access.AuthenticationID, newPass)
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	var checkAuth types.CheckAuth

	err := json.NewDecoder(r.Body).Decode(&checkAuth)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var token *jwt.Token

	token, err = jwt.Parse(checkAuth.Token, func(token *jwt.Token) (interface{}, error) {
		return config.VerifyKey, nil
	})
	if token == nil || err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	if !token.Valid {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	id, ok := token.Claims.(jwt.MapClaims)["jti"].(string)
	userID, _ := strconv.Atoi(id)
	if !ok || userID < 1 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	_, ok = token.Claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	access, err := entity.GetAppUserByID(userID)
	if err != nil {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	claims := jwt.StandardClaims{
		Id:        strconv.Itoa(access.ID), //Todo mudar para jti , id do jtw autenticado.
		Subject:   access.Username,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(config.Config.TokenExpires)).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	////log.Println(claims)
	newToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := newToken.SignedString(config.SignKey)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	response := types.UserResponse{
		AccessToken: tokenString,
		TokenType:   "RS256",
		Username:    access.Username,
	}

	config.JSONResponse(response, http.StatusOK, w)
	return
}
