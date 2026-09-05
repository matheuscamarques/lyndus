package appuser

import (
	"bitbucket.org/lyndus/backend/infra/utils"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"github.com/go-chi/chi"
)

func GetBSsMap(w http.ResponseWriter, r *http.Request) {
	latSTR := r.URL.Query().Get("lat")
	lonSTR := r.URL.Query().Get("lon")
	var err error
	categorySTR := r.URL.Query().Get("category")
	categoryID, _ := strconv.Atoi(categorySTR)

	lat, err := strconv.ParseFloat(latSTR, 64)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	lon, err := strconv.ParseFloat(lonSTR, 64)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	//if err != nil {
	//	errs := []config.Err{
	//		{
	//			Message: "Categoria inválida.",
	//			Code:    "",
	//		},
	//	}
	//	config.JSONResponse(config.ResponseErr{Errs: errs}, api.StatusBadRequest, w)
	//	return
	//}
	var bss []entity.BS
	if categoryID > 0 {
		bss, err = services.BSService.GetBSsByCategoryMap(categoryID, lat, lon)
	} else {
		bss, err = services.BSService.GetBSsMap(lat, lon)
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if len(bss) == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	config.JSONResponse(bss, http.StatusOK, w)
}

func GetBSs(w http.ResponseWriter, r *http.Request) {

	var bss criteria.BSResponse
	var page, itemsPerPage int
	var fantasyName string

	var err error
	categoryID, _ := strconv.Atoi(r.URL.Query().Get("category"))
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ = strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	fantasyName = r.URL.Query().Get("name")

	BitSize := 64
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), BitSize)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), BitSize)

	if page == 0 {
		page = 1
	}

	if categoryID > 0 {
		bss, err = services.BSService.GetBSsByCategory(categoryID, page, itemsPerPage, fantasyName, lat, lon)
	} else {
		bss, err = services.BSService.GetBSs(page, itemsPerPage, fantasyName, lat, lon)
	}

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if len(bss.Items) == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	config.JSONResponse(bss, http.StatusOK, w)
}

func GetBS(w http.ResponseWriter, r *http.Request) {
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bs, err := services.BSService.GetBS(bsID)
	if err == sql.ErrNoRows {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(bs, http.StatusOK, w)
}

func GetBSLogo(w http.ResponseWriter, r *http.Request) {
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	filePath := utils.GetPathFile(bsID, bsID, config.Config.BaseStaticPath, "bs", "logo", "png")
	f, err := os.Open(filePath)

	if os.IsNotExist(err) {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	fi, err := f.Stat()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename=logo.png")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	//stream the body to the client without fully loading it into memory
	_, err = io.Copy(w, f)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
}

func GetBsBMLogo(w http.ResponseWriter, r *http.Request) {
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bmID, err := strconv.Atoi(chi.URLParam(r, "bmID"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bmName, err := services.BSService.GetBSBMName(bsID, bmID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if bmName == "" {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	filePath := utils.GetPathFile(bsID, bmID, config.Config.BaseStaticPath, "bs", "bs_bm", "png")
	f, err := os.Open(filePath)

	if os.IsNotExist(err) {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	fi, err := f.Stat()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename="+bmName+".png")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	//stream the body to the client without fully loading it into memory
	_, err = io.Copy(w, f)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
}

func GetBsBMLogo64(w http.ResponseWriter, r *http.Request) {
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bmID, err := strconv.Atoi(chi.URLParam(r, "bmID"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bmName, err := services.BSService.GetBSBMName(bsID, bmID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if bmName == "" {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	filePath := utils.GetPathFile(bsID, bmID, config.Config.BaseStaticPath, "bs", "bs_bm", "png")

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	image, err := utils.PathImageToBase64Str(filePath)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(image, http.StatusOK, w)

}

func GetBSBanners(w http.ResponseWriter, r *http.Request) {
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	banners, err := services.BSBannerService.GetBSBanners(bsID)
	log.Println(banners)
	if err == sql.ErrNoRows {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	var filePath string
	for i := range banners {
		filePath = utils.GetPathFile(bsID, banners[i].ID, config.Config.BaseStaticPath, "bs", "banner", "png")
		if _, err := os.Stat(filePath); err != nil {
			continue
		}
		banners[i].URL = "https://appuser.lyndus.com.br/api/v1/bs/" + strconv.Itoa(bsID) + "/banners/" + strconv.Itoa(banners[i].ID)
		//banners[i].URL = "http://localhost:8000/api/v1/bs/" + strconv.Itoa(bsID) + "/banners/" + strconv.Itoa(banners[i].ID)
	}

	config.JSONResponse(banners, http.StatusOK, w)
}

func GetBSBanner(w http.ResponseWriter, r *http.Request) {
	log.Println("banana")
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bannerID, err := strconv.Atoi(chi.URLParam(r, "bannerID"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	filePath := utils.GetPathFile(bsID, bannerID, config.Config.BaseStaticPath, "bs", "banner", "png")
	f, err := os.Open(filePath)

	if os.IsNotExist(err) {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	fi, err := f.Stat()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename=logo.png")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	//stream the body to the client without fully loading it into memory
	_, err = io.Copy(w, f)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
}

func GetBsServices(w http.ResponseWriter, r *http.Request) {

	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
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

	var servicesCategories []entity.ServiceCategory
	servicesCategories, err = services.AppUserService.GetServiceCategory()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var remove []int
	for i := range servicesCategories {
		servicesCategories[i].Services, err = services.AppUserService.GetServicesByCategory(bsID, servicesCategories[i].ID)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		if servicesCategories[i].Services == nil {
			remove = append(remove, i)
		}
	}

	for range remove {
		for i := range servicesCategories {
			if servicesCategories[i].Services == nil {
				if i == len(servicesCategories) {
					servicesCategories = servicesCategories[:i]
				} else {
					servicesCategories = append(servicesCategories[:i], servicesCategories[i+1:]...)
				}
				break
			}
		}
	}

	if len(servicesCategories) <= 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	config.JSONResponse(servicesCategories, http.StatusOK, w)
}

func GetBSsWithLocalization(w http.ResponseWriter, r *http.Request) {
	BitSize := 64
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), BitSize)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), BitSize)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	categorySTR := r.URL.Query().Get("category")
	categoryID, _ := strconv.Atoi(categorySTR)

	var bss criteria.BSResponse
	var page, itemsPerPage int
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ = strconv.Atoi(r.URL.Query().Get("itemsPerPage"))

	if categoryID > 0 {
		bss, err = services.BSService.GetBSsByCategoryLocalization(categoryID, lat, lon, page, itemsPerPage)
	} else {
		bss, err = services.BSService.GetBSsWithLocalization(lat, lon, page, itemsPerPage)
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	config.JSONResponse(bss, http.StatusOK, w)
}
