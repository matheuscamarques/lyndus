package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func PostSocial(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var socialRequest entity.BsSocial
	err = json.NewDecoder(r.Body).Decode(&socialRequest)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	socialRequest.BSID = bsID

	social, err := services.BSSocialService.GetSocialByID(socialRequest.SocialID)

	if social.SocialID != socialRequest.SocialID || social.SocialID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	_, err = services.BSSocialService.Create(socialRequest)
	if err != nil {

		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func PutSocial(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var socialRequest entity.BsSocial
	err = json.NewDecoder(r.Body).Decode(&socialRequest)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	socialRequest.BSID = bsID

	social, err := services.BSSocialService.GetBSSocialByID(bsID, socialRequest.SocialID)

	if social.ID == 0 || socialRequest.Url == "" {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BSSocialService.Update(socialRequest)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func DeleteSocial(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var socialRequest entity.BsSocial
	socialRequest.SocialID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	socialRequest.BSID = bsID

	social, err := services.BSSocialService.GetBSSocialByID(bsID, socialRequest.SocialID)
	if social.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BSSocialService.Remove(socialRequest.BSID, socialRequest.SocialID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func GetSocial(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bsSocial, err := services.BSSocialService.GetAll(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(bsSocial, http.StatusOK, w)
}

func ListSocials(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	socials, err := services.BSSocialService.GetSocials(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(socials, http.StatusOK, w)
}
