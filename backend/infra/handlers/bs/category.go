package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.CATEGORY, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}


	categories, err := services.CategoryService.GetBSCategories(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(categories, http.StatusOK, w)
}

func AddCategories(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.CATEGORY, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var catID response.ID
	err = json.NewDecoder(r.Body).Decode(&catID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if catID.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.CategoryService.AddBSCategory(bsID, catID.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	catID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.CategoryService.DeleteBSCategory(bsID, catID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
