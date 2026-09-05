package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

func ProductBrandSave(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var brand entity.ProductBrand
	err = json.NewDecoder(r.Body).Decode(&brand)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var respID response.ID
	respID.ID, err = services.ProductBrandService.Create(brand)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	config.JSONResponse(respID, http.StatusOK, w)
}

func ProductBrandGet(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SUPPLIER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	brandID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	brand, err := services.ProductBrandService.GetByID(brandID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(brand, http.StatusOK, w)
}


func ProductBrandList(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SUPPLIER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	//if err != nil {
	//	config.ResponsePerErr(w, err, config.INVALIDREQUEST)
	//	return
	//}

	itemsPerPag, err := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	if itemsPerPag == 0 {
		itemsPerPag = 10
	}
	//if err != nil {
	//	config.ResponsePerErr(w, err, config.INVALIDREQUEST)
	//	return
	//}

	response, err := services.ProductBrandService.GetAll(page, itemsPerPag)
	if err != nil {
		fmt.Println(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(response, http.StatusOK, w)

}

func ProductBrandSearch(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SUPPLIER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	name := chi.URLParam(r, "name")



	response, err := services.ProductBrandService.SearchByName(name)
	if err != nil {
		fmt.Println(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	log.Println(response, err)
	if response == nil{
		response = make([]entity.ProductBrand, 0)
	}
	log.Println(response, err)

	config.JSONResponse(response, http.StatusOK, w)
}