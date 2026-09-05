package appuser

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {

	categories, err := services.CategoryService.GetCategories()
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBESAVED)
		return
	}
	config.JSONResponse(categories, http.StatusOK, w)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))

	category := entity.Category{
		ID: categoryID,
	}

	err = services.CategoryService.GetCategory(&category)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	bs, err := services.BSService.GetBSsByCategory(categoryID, page, itemsPerPage, "", 0, 0)
	category.BS = bs.Items

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(category, http.StatusOK, w)
}
