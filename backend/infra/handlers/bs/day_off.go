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

func GetBSDayOff(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	daysOff, err := services.DayOffService.GetBSDaysOff(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(daysOff, http.StatusOK, w)
}

func AddBSDayOff(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var dayOff entity.DayOff
	err = json.NewDecoder(r.Body).Decode(&dayOff)
	if err != nil || dayOff.EndDate.IsZero() || dayOff.StartDate.IsZero() {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if dayOff.StartDate.After(dayOff.EndDate.Time){
		config.ResponsePerErr(w, err, config.DATAINVALID)
		return
	}

	var resp response.ID
	resp.ID, err = services.DayOffService.AddBSDayOff(bsID, dayOff)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(resp, http.StatusOK, w)
	return
}

func DeleteBSDayOff(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	dayOffID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.DayOffService.DeleteBSDayOff(bsID, dayOffID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}



func GetBSBMDayOff(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bmID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	daysOff, err := services.DayOffService.GetBSBMDaysOff(bsID, bmID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(daysOff, http.StatusOK, w)
}

func AddBSBMDayOff(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	bmID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var dayOff entity.DayOff
	err = json.NewDecoder(r.Body).Decode(&dayOff)
	if err != nil || dayOff.EndDate.IsZero() || dayOff.StartDate.IsZero() {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if dayOff.StartDate.After(dayOff.EndDate.Time){
		config.ResponsePerErr(w, err, config.DATAINVALID)
		return
	}

	var resp response.ID
	resp.ID, err = services.DayOffService.AddBSBMDayOff(bsID, bmID, dayOff)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(resp, http.StatusOK, w)
	return
}

func DeleteBSBMDayOff(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	bmID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	dayOffID, err := strconv.Atoi(chi.URLParam(r, "day_id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.DayOffService.DeleteBSBMDayOff(bsID, bmID, dayOffID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}


