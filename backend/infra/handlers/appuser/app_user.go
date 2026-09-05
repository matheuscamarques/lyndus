package appuser

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/services"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/types"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"

	"net/http"
	"strconv"
)

const modulePath = "app_user"

func NotRegistered(w http.ResponseWriter, r *http.Request) {
	var user types.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	access, err := entity.GetAppUserUserEmailBYUsername(user.Username)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if access.ID > 0 {
		config.ResponsePerErr(w, nil, config.CPFALREADYREGISTERED)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var appUserCreate entity.AppUserCreate
	err := json.NewDecoder(r.Body).Decode(&appUserCreate)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	genderID, err := services.UserGenderService.GetGenderID(appUserCreate.GenderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if genderID == 0 {
		config.ResponsePerErr(w, err, config.DATAINVALID)
		return
	}

	access, err := entity.GetUserByCredential(appUserCreate.CPF)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if access.ID > 0 {
		config.ResponsePerErr(w, err, config.CPFALREADYREGISTERED)
		return
	}

	id, err := services.CreateService.GetPersonIDByCPF(&appUserCreate)
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBESAVED)
		return
	}

	if id > 0 {
		appUserCreate.ID = id
		err = services.CreateService.UpdatePerson(appUserCreate)
		if err != nil {
			config.ResponsePerErr(w, err, config.CANNOTBESAVED)
			return
		}
	} else {
		appUserCreate.ID, err = services.CreateService.CreatePerson(appUserCreate)
		if err != nil {
			config.ResponsePerErr(w, err, config.CANNOTBESAVED)
			return
		}
	}

	if appUserCreate.ID == 0 {
		config.ResponsePerErr(w, err, config.CANNOTBESAVED)
		return
	}

	authenticator := auth.NewAuthenticationService()
	authID, err := authenticator.CreateAuthentication(appUserCreate.Password)
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBESAVED)
		return
	}

	appID, err := services.CreateService.GetAppUserIDByCPF(&appUserCreate)
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBESAVED)
		return
	}

	appUser := entity.AppUser{
		ID:               appID,
		AuthenticationID: authID,
		PersonID:         appUserCreate.ID,
		Email:            appUserCreate.Email,
		AcceptedTerm:     appUserCreate.AcceptedTerm,
		CPF:              appUserCreate.CPF,
	}
	if appID > 0 {

		err = services.UserUserService.UpdateAppUser(appUser)
	} else {
		id, err = services.UserUserService.CreateAppUser(appUser)
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBESAVED)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func GetGenders(w http.ResponseWriter, r *http.Request) {
	genders, err := services.UserGenderService.GetGenders()
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBESAVED)
		return
	}
	config.JSONResponse(genders, http.StatusOK, w)
}


func GetBsBmsByService(w http.ResponseWriter, r *http.Request) {
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	serviceID, err := strconv.Atoi(chi.URLParam(r, "serviceID"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bsIDdb, err := services.BSService.GetBSID(bsID)
	if err == sql.ErrNoRows || bsIDdb == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	bms, err := services.BM.GetBMsBasicByServiceID(bsID, serviceID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	for k := range bms{
		bms[k].WeekDays, err = services.BM.GetBMWeekDaysWithName(bms[k].ID)
		if err != nil{
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}

	config.JSONResponse(bms, http.StatusOK, w)
}


