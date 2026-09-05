package client

import (
	"bitbucket.org/lyndus/backend/domain/client"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"encoding/json"

	"net/http"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	clientID, err := services.AuthService.CheckPermission(id, client.PROFILE, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}


	company, err := services.CompanyService.GetCompanyByID(clientID)
	if err != nil {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	config.JSONResponse(company, http.StatusOK, w)
}
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var company entity.Company

	id  := r.Context().Value("id").(int)

	clientID, err := services.AuthService.CheckPermission(id, client.PROFILE, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	company.ID = clientID

	err = services.CompanyService.UpdateCompany(company)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(company, http.StatusOK, w)
}
