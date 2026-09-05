package bs

import (
	"encoding/json"

	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/utils"
	"github.com/go-chi/chi"

	"net/http"
	"os"
	"strconv"
)

//SaveBanners func set profile Logo
func SaveBanners(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.BANNER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var banner entity.Banner
	err = json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	banner.ID, err = services.BannerService.CreateBanner(bsID, banner.Name)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	filePath := utils.GetPathFile(bsID, banner.ID, config.Config.BaseStaticPath, modulePath, "banner", "png")

	err = utils.Base64ImageSave(banner.Image, filePath)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(response.ID{ID: banner.ID}, http.StatusOK, w)
}

//GetBanners func get bs banners
func GetBanners(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BANNER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	banners, err := services.BannerService.GetBanners(bsID)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}

	var filePath string
	for k := range banners {
		filePath = utils.GetPathFile(bsID, banners[k].ID, config.Config.BaseStaticPath, modulePath, "banner", "png")
		banners[k].Image, err = utils.PathImageToBase64Str(filePath)
		if err != nil {
			//log.Println(err)
		}
	}
	config.JSONResponse(banners, http.StatusOK, w)
}

func DeleteBanner(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BANNER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bannerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
	}

	err = services.BannerService.DeleteBanner(bsID, bannerID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}

	originPath := utils.GetPathFile(bsID, bannerID, config.Config.BaseStaticPath, modulePath, "banner", "png")
	destinationPath := utils.GetPathFile(bsID, bannerID, config.Config.BaseDeletedPath, modulePath, "banner", "png")

	err = os.Rename(originPath, destinationPath)
	if err != nil {
		//log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	return
}
