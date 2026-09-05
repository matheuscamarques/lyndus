package bonder

import (
	"bitbucket.org/lyndus/backend/domain/bonder"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"

	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/domain/bonder/services"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/go-chi/chi"
)

type ClientRegister struct {
	ID       int       `json:"id"`
	Password string    `json:"password"`
	Company  entity.BS `json:"company"`
}

func ClientRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var clientRegister ClientRegister
	err := json.NewDecoder(r.Body).Decode(&clientRegister)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	clientID, err := services.BonderClient.SearchByCNPJ(clientRegister.Company.CNPJ)

	if clientID != 0 {
		config.ResponsePerErr(w, err, config.CNPJALREADYREGISTERED)
		return
	}

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	authID, err := auth.Service.CreateAuthentication(clientRegister.Password)
	if err != nil {
		config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
		return
	}

	companyID, err := services.BonderCompany.Register(&clientRegister.Company)

	if err != nil && companyID != 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if err != nil {
		config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
		return
	}

	var clientR entity.Client

	clientR.CompanyID = companyID
	clientR.CNPJ = clientRegister.Company.CNPJ

	clientR.ID, err = services.BonderClient.Register(clientR)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	// Criar Categoria Default
	category := entity.DefaultClientCategory(clientR.ID)
	_, err = services.ClientCategory.CreateCategory(category.Category)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	// TODO CRIAR USUARIO MASTER 0001 client.AuthenticationID = authID
	// authID
	clientU := entity.ClientUser{}
	clientU.ClientID = clientR.ID
	clientU.AuthenticationID = authID
	clientU.Master = true
	clientU.CreatedBy = 0 // todo ad
	clientU.Email = clientRegister.Company.Email
	clientU.Name = clientRegister.Company.CompanyName
	clientU.Phone = &clientRegister.Company.Phone

	clientU.ClientUser, err = services.ClientUser.Register(clientU.ClientUser)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	type Response struct {
		ID       int `json:"id,omitempty"`
		CNPJ     types.CNPJ
		Username string
	}

	config.JSONResponse(Response{
		ID:       clientR.ID,
		CNPJ:     clientR.CNPJ,
		Username: clientU.Username,
	}, http.StatusOK, w)
}

func ClientGetAllHandler(w http.ResponseWriter, r *http.Request) {

	var clientResponse response.ListResponse
	var err error

	clientResponse.ActivePage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	search := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("orderBy")
	sortDesc, _ := strconv.ParseBool(r.URL.Query().Get("sortDesc"))
	activeSTR := r.URL.Query().Get("active")
	if activeSTR != "" {
		clientResponse.Active, _ = strconv.ParseBool(activeSTR)
	}

	clientResponse.Items, clientResponse.TotalItems, clientResponse.TotalPages, err = services.BonderClient.GetAll(
		clientResponse.ActivePage,
		itemsPerPage,
		search,
		orderBy,
		sortDesc,
		clientResponse.Active)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(clientResponse, http.StatusOK, w)
}

func ClientGetByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	response, err := services.BonderClient.GetById(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(response, http.StatusOK, w)
}

func ClientActiveHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if id == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BonderClient.UpdateLyndusActive(id, true)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func ClientInactiveHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if id == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BonderClient.UpdateLyndusActive(id, false)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func ClientDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if id == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BonderClient.Detete(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func ClientUpdateHandler(w http.ResponseWriter, r *http.Request) {

	var clientRegister ClientRegister
	err := json.NewDecoder(r.Body).Decode(&clientRegister)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	clientID, err := services.BonderClient.SearchByCNPJ(clientRegister.Company.CNPJ)

	if clientID != clientRegister.ID {
		//fmt.Println(clientID , clientRegister.ID)
		config.ResponsePerErr(w, err, config.CNPJINVALID)
		return
	}

	if err != nil {
		config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
		return
	}

	company := entity.BS{}
	company = clientRegister.Company
	company.ID, err = services.BonderClient.GetCompanyId(clientID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	err = services.BonderCompany.Update(&company)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	// TODO change password
	// puxar o master

	w.WriteHeader(http.StatusOK)
	return
}

func BenefitsGetAllHandler(w http.ResponseWriter, r *http.Request) {

	var respList response.ListResponse

	respList.ActivePage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	search := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("orderBy")
	sortDesc, _ := strconv.ParseBool(r.URL.Query().Get("sortDesc"))

	if respList.ActivePage < 1 {
		respList.ActivePage = 1
	}

	var err error

	respList.Items, respList.TotalItems, respList.TotalPages, err = services.BenefitService.GetBenefits(respList.ActivePage, itemsPerPage, search, orderBy, sortDesc)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(respList, http.StatusOK, w)
}

func BenefitGetHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if id < 1 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	benefit, err := services.BenefitService.GetBenefit(id)
	if err != nil {
		fmt.Println(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(benefit, http.StatusOK, w)
}

func BenefitSendTicketHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if id < 1 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20 Mb
	err = r.ParseMultipartForm(20 * 1024 * 1024)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDMAXSIZE)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			config.ResponsePerErr(w, err, config.FILEINVALID)
			return
		}
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	fileType := http.DetectContentType(buff)
	if fileType != "application/pdf" {
		config.ResponsePerErr(w, err, config.FILEINVALID)
		return
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	benefitID, clientID, statusID, err := services.BenefitService.GetBenefitClientID(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	filePath := utils.GetPathFile(clientID, benefitID, config.Config.BaseStaticPath, "client", "benefit_", "pdf")

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		panic(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	err = r.MultipartForm.RemoveAll()
	if err != nil {
		logger.Error("Erro ao deletar Multi Part Form:", zap.String("error", err.Error()))
	}

	if statusID == bonder.BENEFITSTATUSREQUESTBILLINGTICKET {
		err = services.BenefitService.UpdateBenefitStatus(clientID, benefitID, bonder.BENEFITSTATUSWAITINGPAYMENT)
		if err != nil {
			config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
			return
		}
		config.JSONResponse(response.BenefitStatus{ID: bonder.BENEFITSTATUSWAITINGPAYMENT, Status: bonder.BENEFITSTATUSPT[bonder.BENEFITSTATUSWAITINGPAYMENT]}, http.StatusOK, w)
		return
	}

	config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
	return
}

func BenefitConfirmPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if id < 1 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	//benefitID, clientID, statusID, err := services.BenefitService.GetBenefitClientID(id)
	//if err != nil {
	//	config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	//	return
	//}

	benefit, err := services.BenefitService.GetBenefit(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if benefit.StatusID != bonder.BENEFITSTATUSWAITINGPAYMENT && benefit.StatusID != bonder.BENEFITSTATUSPAYMENTCONFIRMED {
		config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
		return
	}

	benefitUsers, err := services.BenefitService.GetBenefitsUsers(benefit.ID, benefit.ClientID)
	for k := range benefitUsers {
		err = services.BenefitService.PaymentBenefit(
			benefitUsers[k].AppUserID,
			benefit.ClientID,
			benefitUsers[k].Value.Add(benefitUsers[k].AdditionalValue),
			benefit.FantasyName)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}

	err = services.BenefitService.UpdateBenefitStatus(benefit.ClientID, benefit.ID, bonder.BENEFITSTATUSPAID)
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
		return
	}
	config.JSONResponse(response.BenefitStatus{ID: bonder.BENEFITSTATUSPAID, Status: bonder.BENEFITSTATUSPT[bonder.BENEFITSTATUSPAID]}, http.StatusOK, w)
	return
}
