package client

import (
	"bitbucket.org/lyndus/backend/internal/api"
	"errors"
	"net/http"
	"strconv"

	"bitbucket.org/lyndus/backend/domain/client"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/services"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

//{
//	"name"  : "Matheus"
//	"email" : "matheuscamarques@gmail.com"
//	"phone" : "4132661779"
//	permissions:[
//	{"id":1, "level_id":3},
//	{"id":2, "level_id":2},
//	{"id":3, "level_id":1},
//	{"id":4, "level_id":2},
//	],
//	"status": true,
//	"password" : "hash265"
//}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	user := entity.ClientUser{}
	err := types.Object{}.AssignBody(&user, r.Body)

	clientID, err := services.AuthService.CheckPermission(id, client.USER, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	// Cria Uma Autenticação
	authID, err := auth.Service.CreateAuthentication(user.Password)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	// Vincula Autenticação id ao bs_user
	user.AuthenticationID = authID
	user.CreatedBy = id
	user.ClientID = clientID

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

	id := r.Context().Value("id").(int)

	clientID, err := services.AuthService.CheckPermission(id, client.USER, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	user := entity.ClientUser{}
	err = types.Object{}.AssignBody(&user, r.Body)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	//todo checar se usuário existe

	err = services.UserService.Update(user)
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
	authID := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(authID, client.USER, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	param := chi.URLParam(r, "id")
	clientUserID, err := strconv.Atoi(param)
	if err != nil || clientUserID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	user, err := services.UserService.Get(clientID, clientUserID)
	user.Status = config.UserStatusPt[user.StatusID]
	if user.StatusID == config.USERSTATUSACTIVE {
		user.Active = true
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(user, http.StatusOK, w)
	return
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	authID := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(authID, client.USER, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var usersResponse api.ListResponse
	usersResponse.Active = true

	usersResponse.ActivePage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	search := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("orderBy")
	sortDesc, _ := strconv.ParseBool(r.URL.Query().Get("sortDesc"))
	activeSTR := r.URL.Query().Get("active")
	if activeSTR != "" {
		usersResponse.Active, _ = strconv.ParseBool(activeSTR)
	}
	if usersResponse.ActivePage == 0 {
		usersResponse.ActivePage = 1
	}

	usersResponse.Items, usersResponse.TotalItems, usersResponse.TotalPages, err =
		services.UserService.GetAllV2(
			usersResponse.ActivePage,
			itemsPerPage,
			clientID,
			search,
			orderBy,
			sortDesc,
			usersResponse.Active)

	if err != nil {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	config.JSONResponse(usersResponse, http.StatusOK, w)
	return
}

func GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	authID := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(authID, client.USER, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	param := chi.URLParam(r, "id")
	clientUserID, err := strconv.Atoi(param)
	if err != nil || clientUserID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	permissions, err := services.UserService.UserPermissions(clientUserID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(permissions, http.StatusOK, w)
	return
}

func ActiveUser(w http.ResponseWriter, r *http.Request) {
	authID := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(authID, client.USER, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	param := chi.URLParam(r, "id")
	clientUserID, err := strconv.Atoi(param)
	if err != nil || clientUserID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.UserService.UpdateStatus(clientID, clientUserID, config.USERSTATUSACTIVE)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func InactiveUser(w http.ResponseWriter, r *http.Request) {
	authID := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(authID, client.USER, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	param := chi.URLParam(r, "id")
	clientUserID, err := strconv.Atoi(param)
	if err != nil || clientUserID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	user, err := services.UserService.Get(clientID, clientUserID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if user.Username == "001" || user.Master {
		err = errors.New("administrador não pode ser alterado")
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}

	err = services.UserService.UpdateStatus(clientID, clientUserID, config.USERSTATUSINACTIVE)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	authID := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(authID, client.USER, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	param := chi.URLParam(r, "id")
	clientUserID, err := strconv.Atoi(param)
	if err != nil || clientUserID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	user, err := services.UserService.Get(clientID, clientUserID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if user.Username == "001" || user.Master {
		err = errors.New("administrador não pode ser alterado")
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}

	err = services.UserService.UpdateStatus(clientID, clientUserID, config.USERSTATUSDELETED)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
