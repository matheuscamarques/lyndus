package client

import (
	"bitbucket.org/lyndus/backend/internal/api"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"

	"bitbucket.org/lyndus/backend/domain/client"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
	"bitbucket.org/lyndus/backend/domain/client/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/rest/response"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	clientID, err := services.AuthService.CheckPermission(id, client.CATEGORY, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var category entity.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var catID response.ID
	category.ClientID = clientID

	catID.ID, err = services.CategoryService.CreateCategory(category)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	config.JSONResponse(catID, http.StatusOK, w)
}
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.CATEGORY, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var category entity.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	category.ClientID = clientID
	repo := repository.NewClientCategoryRepository()
	err = repo.UpdateCategory(category)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	w.WriteHeader(http.StatusOK)
	return
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.CATEGORY, client.VIEW)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var categoriesResponse api.ListResponse
	categoriesResponse.Active = true

	categoriesResponse.ActivePage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	search := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("orderBy")
	sortDesc, _ := strconv.ParseBool(r.URL.Query().Get("sortDesc"))
	activeSTR := r.URL.Query().Get("active")
	if activeSTR != "" {
		categoriesResponse.Active, _ = strconv.ParseBool(activeSTR)
	}

	categoriesResponse.Items, categoriesResponse.TotalItems, categoriesResponse.TotalPages, err = services.CategoryService.GetCategoriesV2(
		categoriesResponse.ActivePage,
		itemsPerPage,
		clientID,
		search,
		orderBy,
		sortDesc,
		categoriesResponse.Active)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(categoriesResponse, http.StatusOK, w)
}

func GetCategoryList(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.CATEGORY, client.VIEW)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	categories, err := services.CategoryService.GetCategories(clientID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(categories, http.StatusOK, w)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.CATEGORY, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	idCategory, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	category := entity.Category{
		ID:       idCategory,
		ClientID: clientID,
	}

	err = services.CategoryService.GetCategoryByID(&category)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	config.JSONResponse(category, http.StatusOK, w)
	return
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.CATEGORY, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	idCategory, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.CategoryService.DeleteCategory(clientID, idCategory)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}

	w.WriteHeader(http.StatusOK)
	return
}

func ActiveCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.CATEGORY, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	idCategory, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.CategoryService.ActiveInactiveCategory(clientID, idCategory, true)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}

	w.WriteHeader(http.StatusOK)
	return
}

func InactiveCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.CATEGORY, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	idCategory, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.CategoryService.ActiveInactiveCategory(clientID, idCategory, false)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}

	w.WriteHeader(http.StatusOK)
	return
}
