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
	"net/http"
	"strconv"
)

func SupplierList(w http.ResponseWriter, r *http.Request) {
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

	response, err := services.SupplierService.GetAll(bsID, page, itemsPerPag)
	if err != nil {
		fmt.Println(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(response, http.StatusOK, w)

}

//SupplierSave create new bs_supplier
func SupplierSave(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SUPPLIER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var supplier entity.Supplier
	err = json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	supplier.BsID = bsID
	supplier.ID, err = services.SupplierService.Create(supplier)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	resp := response.ID{
		ID: supplier.ID,
	}
	config.JSONResponse(resp, http.StatusOK, w)
}

//SupplierUpdate update bs_supplier
func SupplierUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SUPPLIER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var supplier entity.Supplier
	err = json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	supplier.BsID = bsID

	err = services.SupplierService.Update(supplier)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func SupplierGet(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SUPPLIER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	supplierID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	supplier, err := services.SupplierService.GetByID(bsID, supplierID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(supplier, http.StatusOK, w)
}
