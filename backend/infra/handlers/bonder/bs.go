package bonder

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/domain/bonder/services"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/types"
	"bitbucket.org/lyndus/backend/infra/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const modulePath = "bs"

type Geo struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type BsRequest struct {
	ID         int       `json:"id,omitempty" db:"id"`
	Password   string    `json:"password"`
	Company    entity.BS `json:"company"`
	Geo        Geo       `json:"geo"`
	Categories []int     `json:"categories" db:"categories"`
}

func BsGetCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	response, err := services.BonderBSCategory.GetAll()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if response == nil {
		response = make([]entity.Category, 0)
	}
	config.JSONResponse(response, http.StatusOK, w)
}

func BsRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var bsRegister BsRequest
	err := json.NewDecoder(r.Body).Decode(&bsRegister)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bsID, err := services.BonderBS.SearchByCNPJ(bsRegister.Company.CNPJ)

	if bsID != 0 {
		config.ResponsePerErr(w, err, config.CNPJALREADYREGISTERED)
		return
	}

	if err != nil {
		config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
		return
	}

	//fmt.Printf("%+v\n", bsRegister)

	authID, err := auth.Service.CreateAuthentication(bsRegister.Password)

	if err != nil {
		config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
		return
	}

	companyID, err := services.BonderCompany.Register(&bsRegister.Company)

	if err != nil && companyID != 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if err != nil {
		config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
		return
	}

	var bsR entity.BS

	bsR.CompanyID = companyID
	bsR.CNPJ = bsRegister.Company.CNPJ
	bsR.Lat = &bsRegister.Geo.Latitude
	bsR.Lon = &bsRegister.Geo.Longitude
	bsR.PlanDay = bsRegister.Company.PlanDay
	bsR.PlanValue = bsRegister.Company.PlanValue
	bsR.BalanceDay = bsRegister.Company.BalanceDay
	bsR.RateAnticipation = bsRegister.Company.RateAnticipation
	bsR.RateAnticipationBm = bsRegister.Company.RateAnticipationBm

	bsR.ID, err = services.BonderBS.Register(bsR)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	for k := range bsRegister.Categories {
		err = services.BonderBSCategory.Create(bsR.ID, bsRegister.Categories[k])
		if err != nil {
			logger.Error("Erro ao salvar categoria do BS", zap.String("Err", err.Error()))
		}
	}

	//criar contas iniciais do BS
	err = services.BonderBS.CreateBasicAccounts(bsR.ID, bsR.BalanceDay, bsR.RateAnticipation)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	// TODO CRIAR USUARIO MASTER 0001 client.AuthenticationID = authID
	// authID
	bsU := entity.BsUser{}
	bsU.BsID = bsR.ID
	bsU.AuthenticationID = authID
	bsU.Master = true
	bsU.CreatedBy = 0 // TODO add id do AD BONDER LOGADO
	bsU.Email = bsRegister.Company.Email
	bsU.Name = bsRegister.Company.CompanyName
	bsU.Phone = &bsRegister.Company.Phone

	bsU.BSUser, err = services.BsUser.Register(bsU.BSUser)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	type Response struct {
		ID       int        `json:"id"`
		CNPJ     types.CNPJ `json:"cnpj"`
		Username string     `json:"username"`
	}

	config.JSONResponse(Response{
		ID:       bsR.ID,
		CNPJ:     bsR.CNPJ,
		Username: bsU.Username,
	}, http.StatusOK, w)
}

func BsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var bsRegister BsRequest
	err := json.NewDecoder(r.Body).Decode(&bsRegister)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if bsRegister.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	//fmt.Printf("%+v\n", bsRegister)

	var dbBS entity.BS
	dbBS, err = services.BonderBS.GetByIdSimple(bsRegister.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if bsRegister.ID != dbBS.ID {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bsRegister.Company.ID = dbBS.CompanyID

	err = services.BonderCompany.Update(&bsRegister.Company)
	if err != nil {
		config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
		return
	}

	dbBS.Lat = &bsRegister.Geo.Latitude
	dbBS.Lon = &bsRegister.Geo.Longitude
	dbBS.PlanDay = bsRegister.Company.PlanDay
	dbBS.PlanValue = bsRegister.Company.PlanValue
	dbBS.BalanceDay = bsRegister.Company.BalanceDay
	dbBS.RateAnticipation = bsRegister.Company.RateAnticipation
	dbBS.RateAnticipationBm = bsRegister.Company.RateAnticipationBm

	err = services.BonderBS.Update(dbBS)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	// TODO otimizar, da pra validar de forma mais discreta
	// TODO exemplo infra/handlers/client/employees.go
	var idsCatOK []int
	categories, err := services.BonderBSCategory.GetByBS(bsRegister.ID)
CATEGORY:
	for k := range categories {
		for j := range bsRegister.Categories {
			if bsRegister.Categories[j] == categories[k].ID {
				idsCatOK = append(idsCatOK, categories[k].ID)
				continue CATEGORY
			}
		}
		err = services.BonderBSCategory.Remove(dbBS.ID, categories[k].ID)
		if err != nil {
			logger.Error("Erro ao remover categoria do BS", zap.String("Err", err.Error()))
		}
	}

CATEGORY2:
	for k := range bsRegister.Categories {
		for j := range idsCatOK {
			if bsRegister.Categories[k] == idsCatOK[j] {
				continue CATEGORY2
			}
		}
		err = services.BonderBSCategory.Create(dbBS.ID, bsRegister.Categories[k])
		if err != nil {
			logger.Error("Erro ao salvar categoria do BS", zap.String("Err", err.Error()))
		}
	}

	//TODO CHANGE PASSWORD
	w.WriteHeader(http.StatusOK)
	return
}

func BsGetAllHandler(w http.ResponseWriter, r *http.Request) {

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
	clientResponse.Items, clientResponse.TotalItems, clientResponse.TotalPages, err = services.BonderBS.GetAll(
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

func BsGetByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	response, err := services.BonderBS.GetById(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	response.Categories, err = services.BonderBSCategory.GetByBS(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(response, http.StatusOK, w)
}

func BsActivateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if id == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BonderBS.UpdateLyndusActive(id, true)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func BsInativateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if id == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BonderBS.UpdateLyndusActive(id, false)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func BsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if id == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BonderBS.Detete(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func BsSaveProfileLogo(w http.ResponseWriter, r *http.Request) {
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
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
	if fileType != "image/jpeg" && fileType != "image/png" {
		config.ResponsePerErr(w, err, config.FILEINVALID)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	filePath := utils.GetPathFile(bsID, bsID, config.Config.BaseStaticPath, modulePath, "logo", "png")

	//log.Println(filePath)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = r.MultipartForm.RemoveAll()
	if err != nil {
		logger.Error("Erro ao deletar Multi Part Form:", zap.String("error", err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	return
}

func BsGetProfileLogo(w http.ResponseWriter, r *http.Request) {
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	filePath := utils.GetPathFile(bsID, bsID, config.Config.BaseStaticPath, modulePath, "logo", "png")
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

func BsGetAllPayments(w http.ResponseWriter, r *http.Request) {

	var financialResponse response.ListResponse
	var err error

	sortDesc := true
	financialResponse.ActivePage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	search := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("orderBy")
	sortDescSTR := r.URL.Query().Get("sortDesc")
	activeSTR := r.URL.Query().Get("active")
	if activeSTR != "" {
		financialResponse.Active, _ = strconv.ParseBool(activeSTR)
	}
	if sortDescSTR != "" {
		sortDesc, _ = strconv.ParseBool(sortDescSTR)
	}

	financialResponse.Items, financialResponse.TotalItems, financialResponse.TotalPages, err = services.BonderBS.GetFinancialList(
		financialResponse.ActivePage,
		itemsPerPage,
		search,
		orderBy,
		sortDesc,
		financialResponse.Active)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(financialResponse, http.StatusOK, w)
}

func BsGetByIdPayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	response, err := services.BonderBS.GetBSWithdrawalCash(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(response, http.StatusOK, w)
}

func BSConfirmPayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if id < 1 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bwc, err := services.BonderBS.GetBSWithdrawalCash(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if bwc.Paid != nil && *bwc.Paid {
		config.ResponsePerErr(w, errors.New("inconsistência nos valores"), config.CANNOTBECHANGED)
		return
	}

	wcbs, err := services.BonderBS.GetBSWithdrawalCashBalance(bwc.BSID, bwc.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var total decimal.Decimal
	for k := range wcbs {
		total = total.Add(wcbs[k].Value)
	}

	div := decimal.NewFromInt(100)
	total = total.Sub(bwc.ValueAnticipation.Mul(bwc.RateAnticipation).Div(div)).Round(2)

	if total.Div(bwc.TotalValue).Round(2).GreaterThan(decimal.NewFromInt(1)) {
		config.ResponsePerErr(w, errors.New("inconsistência nos valores"), config.CANNOTBECHANGED)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20 Mb
	err = r.ParseMultipartForm(20 * 1024 * 1024)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDMAXSIZE)
		return
	}

	file, fh, err := r.FormFile("file")
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
	filename := fh.Filename

	//do a regex in file name to check if is a pdf
	pattern := `(?i)\.pdf$`
	matched, err := regexp.MatchString(pattern, filename)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if fileType != "application/pdf" && !matched {
		config.ResponsePerErr(w, err, config.FILEINVALID)
		return
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	filePath := utils.GetPathFile(bwc.BSID, bwc.ID, config.Config.BaseStaticPath, "bs", "withdrawal_cash_", "pdf")
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
		logger.Warning("Erro ao deletar Multi Part Form:", zap.String("error", err.Error()))
	}

	lyndusValue := bwc.ValueReleased.Add(bwc.ValueAnticipation)
	lyndusValueRate := lyndusValue.Sub(bwc.TotalValue)

	err = services.BonderBS.BSWithdrawalCashConfirm(bwc.BSID, bwc.AccountMovementID, bwc.AccountReceiptID, lyndusValue, lyndusValueRate, wcbs)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func BSPaymentConfirmFile(w http.ResponseWriter, r *http.Request) {

	paymentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bwc, err := services.BonderBS.GetBSWithdrawalCash(paymentID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if bwc.ID < 1 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	filePath := utils.GetPathFile(bwc.BSID, bwc.ID, config.Config.BaseStaticPath, "bs", "withdrawal_cash_", "pdf")
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
	w.Header().Set("Content-Disposition", "attachment; filename=withdrawal_cash_"+strconv.Itoa(bwc.ID)+".pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	//stream the body to the client without fully loading it into memory
	_, err = io.Copy(w, f)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
}
