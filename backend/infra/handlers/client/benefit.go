package client

import (
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/utils"
	"bitbucket.org/lyndus/backend/internal/api"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"

	"bitbucket.org/lyndus/backend/domain/client"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"github.com/go-chi/chi"
)

func GetBenefits(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.VIEW)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	var benefitResponse api.ListResponse
	benefitResponse.Active = true

	benefitResponse.ActivePage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	search := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("orderBy")
	sortDesc, _ := strconv.ParseBool(r.URL.Query().Get("sortDesc"))
	status, _ := strconv.Atoi(r.URL.Query().Get("status"))

	benefitResponse.Items, benefitResponse.TotalItems, benefitResponse.TotalPages, err =
		services.BenefitService.GetBenefitsV2(
			benefitResponse.ActivePage,
			itemsPerPage,
			clientID,
			status,
			search,
			orderBy,
			sortDesc)

	if err != nil {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	config.JSONResponse(benefitResponse, http.StatusOK, w)
}

func GetBenefit(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.VIEW)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var benefit entity.Benefit

	benefitID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	benefit.ID = benefitID
	benefit.ClientID = clientID
	respBenefit, err := services.BenefitService.GetBenefit(benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	respBenefit.BenefitUsers, err = services.BenefitService.GetBenefitsUsers(benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	for k := range respBenefit.BenefitUsers {
		categories, err := services.BenefitService.GetBenefitsUserCategory(respBenefit.BenefitUsers[k])
		if err != nil {
			logger.Error("err list benfit bs_user category", zap.String("err:", err.Error()))
		}

		respBenefit.BenefitUsers[k].Categories = categories

	}

	config.JSONResponse(respBenefit, http.StatusOK, w)
}

func UpdateBenefits(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.ALL)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	benefitID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {

		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var benefitUserUpdate entity.BenefitUserUpdate
	err = json.NewDecoder(r.Body).Decode(&benefitUserUpdate)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	benefit := entity.Benefit{
		ClientID: clientID,
		ID:       benefitID,
	}

	status, err := services.BenefitService.GetBenefitStatus(&benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBESAVED)
		return
	}
	if status != client.BENEFITSTATUSOPEN {
		config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
		return
	}

	benefitUser := entity.BenefitUser{
		ClientID:         clientID,
		BenefitID:        benefitID,
		ID:               benefitUserUpdate.ID,
		AdditionalValue:  &benefitUserUpdate.Value,
		AdditionalReason: &benefitUserUpdate.Reason,
	}
	benefitUser.AdditionalValue = &benefitUserUpdate.Value

	err = services.BenefitService.UpdateBenefitUserAdditional(&benefitUser)
	if err != nil {
		config.ResponsePerErr(w, err, config.CANNOTBESAVED)
		return
	}

	benefitsUser, err := services.BenefitService.GetBenefitsUsersValues(benefitUser)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	value := decimal.Zero
	for _, be := range benefitsUser {
		value = value.Add(be.Value)
		if be.AdditionalValue != nil {
			value = value.Add(*be.AdditionalValue)
		}

	}
	benefit.Value = value
	err = services.BenefitService.UpdateBenefitValue(&benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}

	config.JSONResponse(response.Value{Value: value}, http.StatusOK, w)
}
func CreateBenefit(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.ALL)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var benefitCreate entity.BenefitCreate
	err = json.NewDecoder(r.Body).Decode(&benefitCreate)

	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var employeesBenefit []entity.EmployeeBenefit

	if benefitCreate.AllEmployees {
		employeesBenefit, err = services.EmployeeService.GetEmployeesCategoryValue(clientID, true)
	} else if len(benefitCreate.Categories) > 0 {
		employeesBenefit, err = services.EmployeeService.GetEmployeesByCategoryIDValue(clientID, benefitCreate.Categories, true)
	} else {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	benefit := entity.Benefit{
		ClientID:        clientID,
		BenefitStatusID: client.BENEFITSTATUSOPEN,
		Value:           decimal.Zero,
		Desc:            benefitCreate.Desc,
	}

	err = services.BenefitService.CreateBenefit(&benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	//fmt.Println("beneficius", benefit, employeesBenefit)
	type BenefitEmployeeControl struct {
		BenefitUserID int
		Value         decimal.Decimal
	}

	var benefitEmployees map[int]BenefitEmployeeControl
	benefitEmployees = make(map[int]BenefitEmployeeControl)

	benefitUser := entity.BenefitUser{
		ClientID: clientID,
	}
	var buc entity.BenefitUserCategory

	for _, employeeBenefit := range employeesBenefit {
		if value, ok := benefitEmployees[employeeBenefit.ID]; ok {
			value.Value = value.Value.Add(employeeBenefit.CategoryValue)

			benefitUser.BenefitID = value.BenefitUserID
			benefitUser.Value = value.Value

			err = services.BenefitService.UpdateBenefitUserValue(&benefitUser)
			if err != nil {
				panic(err)
			}

			buc.BenefitUserID = value.BenefitUserID
			buc.CategoryID = employeeBenefit.CategoryID
			buc.Name = employeeBenefit.CategoryName
			buc.Value = employeeBenefit.CategoryValue

			benefitEmployees[employeeBenefit.ID] = value
		} else {
			benefitUser.BenefitID = benefit.ID
			benefitUser.AppUserID = employeeBenefit.AppUserID
			benefitUser.Value = employeeBenefit.CategoryValue
			benefitUser.EmployeeID = employeeBenefit.ID

			benefitUser.ID, err = services.BenefitService.CreateBenefitUser(benefitUser)
			if err != nil {
				panic(err)
			}

			buc.BenefitUserID = benefitUser.ID
			buc.CategoryID = employeeBenefit.CategoryID
			buc.Name = employeeBenefit.CategoryName
			buc.Value = employeeBenefit.CategoryValue

			benefitEmployees[employeeBenefit.ID] = BenefitEmployeeControl{
				BenefitUserID: benefitUser.ID,
				Value:         employeeBenefit.CategoryValue,
			}
		}

		//fmt.Println(buc.BenefitUserID)
		//if buc.BenefitUserID == 0 {
		//	continue
		//}
		_, err = services.BenefitService.CreateBenefitUserCategory(&buc)
		if err != nil {
			panic(err)
		}
	}

	//value := decimal.Decimal(0)
	for _, be := range benefitEmployees {
		benefit.Value = benefit.Value.Add(be.Value)
	}

	//fmt.Println("FOI ? ", benefit, benefitEmployees)
	err = services.BenefitService.UpdateBenefitValue(&benefit)
	if err != nil {
		panic(err)
	}

	config.JSONResponse(response.ID{ID: benefit.ID}, http.StatusOK, w)
}

func UpdateBenefit(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.ALL)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var benefit entity.Benefit
	err = json.NewDecoder(r.Body).Decode(&benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if benefit.ID == 0 || benefit.Desc == "" {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	benefit.ClientID = clientID

	status, err := services.BenefitService.GetBenefitStatus(&benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if status != client.BENEFITSTATUSOPEN {
		config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
		return
	}

	err = services.BenefitService.UpdateBenefitDesc(benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func BenefitCancel(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.ALL)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var benefitID response.ID
	err = json.NewDecoder(r.Body).Decode(&benefitID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if benefitID.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	benefit := entity.Benefit{
		ID:       benefitID.ID,
		ClientID: clientID,
	}

	status, err := services.BenefitService.GetBenefitStatus(&benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if status == client.BENEFITSTATUSOPEN ||
		status == client.BENEFITSTATUSWAITINGPAYMENT ||
		status == client.BENEFITSTATUSREQUESTBILLINGTICKET {
		benefit.BenefitStatusID = client.BENEFITSTATUSCANCELED
		err = services.BenefitService.UpdateBenefitStatus(benefit)
		if err != nil {
			config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
			return
		}
		config.JSONResponse(response.BenefitStatus{ID: client.BENEFITSTATUSCANCELED, Status: client.BENEFITSTATUSPT[client.BENEFITSTATUSCANCELED]}, http.StatusOK, w)
		return
	}

	config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
	return
}

func BenefitRequestTicket(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.ALL)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var benefitID response.ID
	err = json.NewDecoder(r.Body).Decode(&benefitID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if benefitID.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	benefit := entity.Benefit{
		ID:       benefitID.ID,
		ClientID: clientID,
	}

	status, err := services.BenefitService.GetBenefitStatus(&benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if status == client.BENEFITSTATUSOPEN {
		benefit.BenefitStatusID = client.BENEFITSTATUSREQUESTBILLINGTICKET
		err = services.BenefitService.UpdateBenefitStatus(benefit)
		if err != nil {
			config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
			return
		}
		config.JSONResponse(response.BenefitStatus{ID: client.BENEFITSTATUSREQUESTBILLINGTICKET, Status: client.BENEFITSTATUSPT[client.BENEFITSTATUSREQUESTBILLINGTICKET]}, http.StatusOK, w)
		return
	}

	config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
	return
}

func BenefitDownloadTicket(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.VIEW)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var benefit entity.Benefit

	benefitID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	benefit.ID = benefitID
	benefit.ClientID = clientID
	respBenefit, err := services.BenefitService.GetBenefit(benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	filePath := utils.GetPathFile(clientID, respBenefit.ID, config.Config.BaseStaticPath, "client", "benefit_", "pdf")
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
	w.Header().Set("Content-Disposition", "attachment; filename="+respBenefit.Desc+".pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	//stream the body to the client without fully loading it into memory
	_, err = io.Copy(w, f)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
}

func BenefitConfirmPayment(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.ALL)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var benefitID response.ID
	err = json.NewDecoder(r.Body).Decode(&benefitID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if benefitID.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	benefit := entity.Benefit{
		ID:       benefitID.ID,
		ClientID: clientID,
	}

	status, err := services.BenefitService.GetBenefitStatus(&benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if status == client.BENEFITSTATUSWAITINGPAYMENT {
		benefit.BenefitStatusID = client.BENEFITSTATUSPAYMENTCONFIRMED
		err = services.BenefitService.UpdateBenefitStatus(benefit)
		if err != nil {
			config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
			return
		}
		config.JSONResponse(response.BenefitStatus{ID: client.BENEFITSTATUSPAYMENTCONFIRMED, Status: client.BENEFITSTATUSPT[client.BENEFITSTATUSPAYMENTCONFIRMED]}, http.StatusOK, w)
		return
	}

	config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
	return
}

func BenefitPayment(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.BENEFIT, client.ALL)
	if err != nil || clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var benefitID response.ID
	err = json.NewDecoder(r.Body).Decode(&benefitID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if benefitID.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	benefit := entity.Benefit{
		ID:       benefitID.ID,
		ClientID: clientID,
	}

	status, err := services.BenefitService.GetBenefitStatus(&benefit)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if status == client.BENEFITSTATUSOPEN {
		benefit.BenefitStatusID = client.BENEFITSTATUSWAITINGPAYMENT
		err = services.BenefitService.UpdateBenefitStatus(benefit)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
	return
}
