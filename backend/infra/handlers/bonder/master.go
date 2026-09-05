package bonder

import (
	"bitbucket.org/lyndus/backend/domain/bonder/services"
	serviceBS "bitbucket.org/lyndus/backend/domain/bs/services"
	serviceClient "bitbucket.org/lyndus/backend/domain/client/services"
	"bitbucket.org/lyndus/backend/global/contracts"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/types"
	"bitbucket.org/lyndus/backend/infra/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func MasterUserBSChangePassword(w http.ResponseWriter, r *http.Request){
	var changePassword types.ChangePassword
	err := json.NewDecoder(r.Body).Decode(&changePassword)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var authentication auth.Authentication
	var service contracts.AuthServiceInterface
	service = serviceBS.NewBSAuthService()
	userID, err := services.BsUser.GetMasterUserAuthID(changePassword.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	authentication, err = service.GetUserPasswordByID(userID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if authentication.ID == 0 {
		config.ResponsePerErr(w, err, config.USERNOTFOUND)
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

func MasterUserClientChangePassword(w http.ResponseWriter, r *http.Request){
	var changePassword types.ChangePassword
	err := json.NewDecoder(r.Body).Decode(&changePassword)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	userID, err := services.ClientUser.GetMasterUserAuthID(changePassword.ID)
	fmt.Println(userID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var authentication auth.Authentication
	var service contracts.AuthServiceInterface = serviceClient.NewClientAuthService()
	authentication, err = service.GetUserPasswordByID(userID)
	fmt.Printf("%+v",authentication)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if authentication.ID == 0 {
		config.ResponsePerErr(w, err, config.USERNOTFOUND)
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