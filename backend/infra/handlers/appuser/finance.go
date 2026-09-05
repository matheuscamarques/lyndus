package appuser

import (
	"bitbucket.org/lyndus/backend/domain/appuser"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/services"
	"bitbucket.org/lyndus/backend/domain/constants"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/utils"
	"bitbucket.org/lyndus/backend/internal/pkg/ebanx"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func GetExtract(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	statements, err := services.StatementService.GetStatements(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if statements == nil {
		statements = make([]entity.Statement, 0)
	}

	//todo pegar fuso horário do app bs_user.
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err == nil {
		for k := range statements {
			statements[k].DateTime.Time = statements[k].DateTime.In(loc)
		}
	}

	config.JSONResponse(statements, http.StatusOK, w)
	return
}

func GetPayments(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	payments, err := services.PaymentService.GetOpenPayments(id, constants.PaymentStatusWaitingPayment)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if payments == nil {
		payments = make([]entity.Payment, 0)
	}

	for k := range payments {

		payments[k].Order, err = services.PaymentService.GetOrderByID(payments[k].BSOrderID)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		payments[k].BsID = &payments[k].Order.BsID
		payments[k].Order.Items, err = services.PaymentService.GetOrderItems(payments[k].BSOrderID)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		payments[k].Order.Payments, err = services.PaymentService.GetOrderPayments(payments[k].BSOrderID)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		payments[k].BSName, err = services.PaymentService.GetBsName(payments[k].Order.BsID)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}

	config.JSONResponse(payments, http.StatusOK, w)
	return
}

func GetBalance(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	balance, err := services.BalanceService.GetBalance(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(balance, http.StatusOK, w)
	return
}

func ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	var confirmPayment entity.ConfirmPayment
	err := json.NewDecoder(r.Body).Decode(&confirmPayment)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var access auth.AuthenticationAppUser
	access, err = entity.GetAppUserUserPasswordByID(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if !utils.ComparePasswords(access.Password, confirmPayment.Password) {
		config.ResponsePerErr(w, err, config.PAYMENTPASSWORD)
		return
	}

	// get bs_payment by uuid
	dbConfirmPayment, err := services.PaymentService.GetPayment(confirmPayment.PaymentUUID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if confirmPayment.Value != dbConfirmPayment.Value {
		config.ResponsePerErr(w, err, config.WRONGVALUE)
		return
	}

	// get balance
	balance, err := services.BalanceService.GetBalance(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if balance.Value < confirmPayment.Value {
		config.ResponsePerErr(w, err, config.BALANCEINSUFFICIENT)
		return
	}

	orderID, err := services.PaymentService.GetOrderID(id, dbConfirmPayment.BsOrderPaymentID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	// update bs_payment status
	err = services.PaymentService.UpdatePayment(dbConfirmPayment.ID, dbConfirmPayment.BsID, constants.PaymentStatusPaid)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	// calculate new balance value
	newValue := balance.Value - dbConfirmPayment.Value

	// update balance value
	err = services.BalanceService.UpdateBalance(id, newValue)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	// insert balance history
	_, err = services.BalanceService.InsertBalanceHistory(id, balance.Value, newValue)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//todo verificar datetime.
	statements := entity.AppUserStatement{
		AppUserID:        id,
		AppUserBalanceID: balance.ID,
		StatementID:      appuser.StatementPayment,
		Value:            dbConfirmPayment.Value,
		Desc:             appuser.Statements[appuser.StatementPayment],
	}
	// insert app bs_user "bank" statement
	_, err = services.PaymentService.InsertStatement(statements)

	// update balance value
	err = services.PaymentService.UpdateBsOrderStatus(orderID, dbConfirmPayment.BsID, constants.OrderStatusPaid)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = services.PaymentService.UpdateScheduleStatus(orderID, dbConfirmPayment.BsID, constants.SCHEDULEFINISHED)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func AddCredtCard(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	var card ebanx.Car
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	tokenizedCard, err := ebanx.Ebanx.CardTokenizer("creditcard", card)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	//todo controle do cartão principal
	// o primeiro cartão se torna o principal
	// demais rotas
	log.Println(tokenizedCard, err)
	err = services.PaymentService.SaveCredCard(id, tokenizedCard.PaymentTypeCode, tokenizedCard.Token, tokenizedCard.MaskedCardNumber)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
