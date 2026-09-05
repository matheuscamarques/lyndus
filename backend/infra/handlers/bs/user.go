package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"log"

	"net/http"
	"strconv"

	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.USER, bs.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	user := entity.BSUser{}
	err = types.Object{}.AssignBody(&user, r.Body)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	// Cria Uma Autenticação
	user.AuthenticationID, err = auth.Service.CreateAuthentication(user.Password)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	user.CreatedBy = id
	user.BsID = bsID
	user, err = services.UserService.Register(user)

	if err != nil {
		logger.Error(
			"Erro ao criar usuário: ",
			zap.String("erro", err.Error()),
		)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	jsonRES := response.ID{
		ID: user.ID,
	}
	config.JSONResponse(jsonRES, http.StatusOK, w)
	return
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var service contracts.BSUserServiceInterface = services.NewBSUserService()
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.USER, bs.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	user := entity.BSUser{}
	err = types.Object{}.AssignBody(&user, r.Body)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	//todo checar se usuario existe e status é ativo

	err = service.Update(user)
	if err != nil {
		logger.Error(
			"Erro ao alterar usuário: ",
			zap.String("erro", err.Error()),
		)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var service contracts.BSUserServiceInterface = services.NewBSUserService()
	authID := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(authID, bs.USER, bs.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	param := chi.URLParam(r, "id")
	BSUserID, err := strconv.Atoi(param)
	if err != nil || BSUserID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	user, err := service.Get(bsID, BSUserID)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(user, http.StatusOK, w)
	return
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {

	authID := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(authID, bs.USER, bs.VIEW)
	log.Println(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	users, err := services.UserService.GetAll(bsID, config.USERSTATUSACTIVE)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(users, http.StatusOK, w)
	return
}
