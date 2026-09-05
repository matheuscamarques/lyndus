package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func GetServiceCategories(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SERVICE, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	serviceCategories, err := services.BS.GetBSServiceCategories()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(serviceCategories, http.StatusOK, w)
}

//SaveService create new bs_service
func SaveService(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.SERVICE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var service entity.Service
	err = json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	service.BsID = bsID

	//repo := repository.NewServiceRepository()

	service.ID, err = services.BS.CreateService(service)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	resp := response.ID{
		ID: service.ID,
	}

	config.JSONResponse(resp, http.StatusOK, w)
}

//UpdateService update bs_service
func UpdateService(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SERVICE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var service entity.Service
	err = json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	service.BsID = bsID

	err = services.BS.UpdateService(service)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func GetServices(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SERVICE, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	name := r.URL.Query().Get("name")
	var serviceList []entity.Service
	if name != "" {
		serviceList, err = services.BS.GetBSServicesByName(bsID, name)
	} else {
		serviceList, err = services.BS.GetBSServices(bsID)
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(serviceList, http.StatusOK, w)
}

func GetService(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	serviceID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || serviceID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	bsID, err := services.AuthService.CheckPermission(id, bs.SERVICE, bs.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	bs, err := services.BS.GetService(bsID, serviceID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(bs, http.StatusOK, w)
}
