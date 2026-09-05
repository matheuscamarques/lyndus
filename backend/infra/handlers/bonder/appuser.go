package bonder

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"

	"bitbucket.org/lyndus/backend/domain/bonder/services"
	"bitbucket.org/lyndus/backend/infra/config"
)

func AppuserGetAll(w http.ResponseWriter, r *http.Request) {

	var appuserResponse response.ListResponse
	var err error

	appuserResponse.ActivePage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	search := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("orderBy")
	sortDesc, _ := strconv.ParseBool(r.URL.Query().Get("sortDesc"))
	activeSTR := r.URL.Query().Get("active")
	if activeSTR != "" {
		appuserResponse.Active, _ = strconv.ParseBool(activeSTR)
	}

	appuserResponse.Items, appuserResponse.TotalItems, appuserResponse.TotalPages, err = services.AppUser.GetAllAppUser(
		appuserResponse.ActivePage,
		itemsPerPage,
		search,
		orderBy,
		sortDesc,
		appuserResponse.Active)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if appuserResponse.Items == nil {
		appuserResponse.Items = make([]entity.AppUserPerson, 0)
	}

	config.JSONResponse(appuserResponse, http.StatusOK, w)
	return
}

func AppUserGetByID(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	response, err := services.AppUser.GetAppUserByID(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	config.JSONResponse(response, http.StatusOK, w)
	return
}

//func AppUserUpdateBonus(w http.ResponseWriter, r *http.Request) {
//	var bonus entity.Bonus
//
//	err := json.NewDecoder(r.Body).Decode(&bonus)
//	if err != nil {
//		config.ResponsePerErr(w, nil, config.INVALIDREQUEST)
//		return
//	}
//
//	err = services.AppUser.UpdateBonus(bonus)
//	log.Println(err)
//	if err != nil {
//		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	return
//}

func AppUserUpdateLyndusBoxBonus(w http.ResponseWriter, r *http.Request) {
	var bonus entity.Bonus

	err := json.NewDecoder(r.Body).Decode(&bonus)
	if err != nil {
		config.ResponsePerErr(w, nil, config.INVALIDREQUEST)
		return
	}

	if bonus.Balance.LessThan(decimal.NewFromInt(1)) {
		config.ResponsePerErr(w, nil, config.INVALIDREQUEST)
		return
	}

	err = services.AppUser.UpdateLyndusBox(bonus)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
