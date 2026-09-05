package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
)

type ProductCaseStockSave struct {
	ID           int             `json:"id"`
	Quantity     int             `json:"quantity"`
	CostPrice    decimal.Decimal `json:"costPrice"`
	SalesPrice   decimal.Decimal `json:"salesPrice"`
	ProfitMargin decimal.Decimal `json:"profitMargin"`
}

func ProductCategories(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	categories, err := services.ProductService.GetCategories()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(categories, http.StatusOK, w)
}

func ProductPackages(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	packages, err := services.ProductService.GetPackages()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(packages, http.StatusOK, w)
}

func ProductUnits(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	units, err := services.ProductService.GetUnits()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(units, http.StatusOK, w)
}

func ProductParameters(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	resp, err := services.ProductService.Parameters()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(resp, http.StatusOK, w)
}

func ProductSave(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var product entity.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if product.ProductClassID == 0 ||
		product.ProductUnitID == 0 ||
		product.ProductPackageID == 0 ||
		product.ProductCategoryID == 0 ||
		product.SupplierID == 0 ||
		product.Quantity == 0 ||
		product.SalesPrice.IsZero() ||
		product.CostPrice.IsNegative() ||
		product.Description == "" {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	//STATUS DO PRODUTO, CONTROLE.
	product.BSID = bsID
	product.CreatedBy = id
	product.Status = bs.ProductStatusActive

	if product.ProductBrandID == 0 {
		if product.BrandName == "" {
			config.ResponsePerErr(w, err, config.INVALIDREQUEST)
			return
		}
		brand, err := services.ProductBrandService.GetByName(product.BrandName)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		if brand.ID == 0 {
			brand.Name = product.BrandName
			brand.ID, err = services.ProductBrandService.Create(brand)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}
		product.ProductBrandID = brand.ID
	}

	n, err := services.ProductService.CountProducts(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//todo gerar código do produto.
	product.Code = fmt.Sprintf("%04d", n+1)

	var respID response.ID
	respID.ID, err = services.ProductService.Create(product)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(respID, http.StatusOK, w)
}

func ProductGet(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	productID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	product, err := services.ProductService.GetByID(bsID, productID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if product.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	//bs_product.SalesPrice = (bs_product.CostPrice / (10000 - bs_product.ProfitMargin)) * bs_product.CostPrice

	config.JSONResponse(product, http.StatusOK, w)
}

func ProductList(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var page, itemsPerPage int
	var productResponse criteria.CProductResponse

	productResponse, err = services.ProductService.GetAllProducts(bsID, page, itemsPerPage)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(productResponse, http.StatusOK, w)
}

func ProductListSimple(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var productResponse []entity.ProductSimple
	productResponse, err = services.ProductService.GetProductsSimple(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	//for k := range productResponse {
	//	productResponse[k].Value = ((productResponse[k].CostPrice * productResponse[k].ProfitMargin) / 10000) + productResponse[k].CostPrice
	//}
	config.JSONResponse(productResponse, http.StatusOK, w)
}

func ProductStockSave(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var product ProductCaseStockSave
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	product.ID, err = strconv.Atoi(chi.URLParam(r, "id"))

	// Atualiza quantidade e valor
	err = services.ProductService.StockSave(bsID, id, product.ID, product.Quantity, product.CostPrice, product.SalesPrice, product.ProfitMargin)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func ProductUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var product entity.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if product.ProductClassID == 0 ||
		product.ProductUnitID == 0 ||
		product.ProductPackageID == 0 ||
		product.ProductCategoryID == 0 ||
		product.SupplierID == 0 ||
		product.Quantity == 0 ||
		product.Description == "" {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	//STATUS DO PRODUTO, CONTROLE.
	product.BSID = bsID
	product.CreatedBy = id
	product.Status = bs.ProductStatusActive

	if product.ProductBrandID == 0 {
		if product.BrandName == "" {
			config.ResponsePerErr(w, err, config.INVALIDREQUEST)
			return
		}
		brand, err := services.ProductBrandService.GetByName(product.BrandName)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		if brand.ID == 0 {
			brand.Name = product.BrandName
			brand.ID, err = services.ProductBrandService.Create(brand)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}
		product.ProductBrandID = brand.ID
	}

	err = services.ProductService.Update(product)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func ProductBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := services.ProductBrandService.GetAll(1, 100)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(brands, http.StatusOK, w)
}

func ProductGetClasses(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PRODUCT, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	config.JSONResponse(bs.ProductClasses, http.StatusOK, w)
}
