package appuser

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/config"
	"net/http"
)

func GetBanners(w http.ResponseWriter, r *http.Request) {
	images, err := entity.Banner{}.GetAppBannerImages()
	if err != nil {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	config.JSONResponse(images, http.StatusOK, w)
}