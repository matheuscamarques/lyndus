package bs

import (
	"bitbucket.org/lyndus/backend/domain/constants"
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"strconv"

	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

//GetPaymentsType LIST payments types
func GetPaymentsType(w http.ResponseWriter, _ *http.Request) {

	payments, err := services.PaymentService.GetPaymentTypes()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(payments, http.StatusOK, w)
}

//SaveOrder create new order
func SaveOrder(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.ORDER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var order entity.Order
	err = json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	personDB, err := services.PersonService.GetPersonID(bsID, order.BSPersonID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if personDB.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	order.BsID = bsID
	order.StatusID = constants.OrderStatusOpen

	order.ID, err = services.OrderService.CreateOrder(order)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	resp := response.ID{
		ID: order.ID,
	}

	config.JSONResponse(resp, http.StatusOK, w)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.ORDER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	orders, err := services.OrderService.GetOrders(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(orders, http.StatusOK, w)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.ORDER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var order entity.Order

	order, err = services.OrderService.GetOrderByID(bsID, orderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if order.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	order.Items, err = services.OrderService.GetOrdersItemsFull(order.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(order, http.StatusOK, w)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.ORDER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var orderDB entity.Order

	orderDB, err = services.OrderService.GetOrderByID(bsID, orderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if orderDB.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	err = services.OrderService.UpdateOrderStatus(bsID, orderID, constants.OrderStatusCanceled)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

//todo transactions on orders and payments

func OrdersCharge(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.ORDER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var paymentsCharge []entity.PaymentCharge
	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&paymentsCharge)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var order entity.Order

	order, err = services.OrderService.GetOrderByID(bsID, orderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if order.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	if order.StatusID == constants.OrderStatusCanceled {
		config.ResponsePerErr(w, err, config.ORDERCANCELED)
		return
	} else if order.StatusID == constants.OrderStatusExpired {
		config.ResponsePerErr(w, err, config.ORDEREXPIRED)
		return
	}

	var value decimal.Decimal
	for _, payment := range paymentsCharge {
		value = value.Add(payment.Value)
	}

	////log.Println(value, order.TotalValue)
	if value != order.TotalValue {
		config.ResponsePerErr(w, err, config.ORDERINVALIDVALUE)
		return
	}

	lyndusPay := false
	for _, payment := range paymentsCharge {
		payment.BSOrderID = orderID
		payment.ID, err = services.PaymentService.CreateOrderPayment(payment)
		if payment.PaymentTypeID == bs.LyndusPayment {
			lyndusPay = true
			// todo pagamento pelo app
			paymentUUID, err := uuid.NewUUID()
			if err != nil {
				logger.Error("Err on generate UUID PAYMENT:", zap.String("ERROR:", err.Error()))
			}
			balanceReceivable := entity.BalanceReceivable{
				BsID:             id,
				BSPersonID:       order.BSPersonID,
				BSOrderPaymentID: payment.ID,
				PaymentStatusID:  bs.PaymentStatusWaitingPayment,
				PaymentUUID:      paymentUUID,
				Value:            payment.Value,
			}
			_, err = services.PaymentService.CreateOrderReceivable(balanceReceivable)
			if err != nil {
				logger.Error("Err on generate UUID PAYMENT:", zap.String("ERROR:", err.Error()))
			}
		}
	}
	if lyndusPay {
		order.StatusID = constants.OrderStatusWaitingPayment
	} else {
		order.StatusID = constants.OrderStatusPaid
	}

	err = services.OrderService.UpdateOrderStatus(bsID, order.ID, order.StatusID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func GetOrdersItems(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.ORDER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var orderDB entity.Order

	orderDB, err = services.OrderService.GetOrderIDByID(bsID, orderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if orderDB.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	orderDB.Items, err = services.OrderService.GetOrdersItemsFull(orderDB.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(orderDB.Items, http.StatusOK, w)
}

func OrdersAddItem(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.ORDER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	var item entity.Item
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var orderDB entity.Order
	orderDB, err = services.OrderService.GetOrderIDByID(bsID, orderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if orderDB.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	if orderDB.StatusID != constants.OrderStatusOpen {
		config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
		return
	}

	if item.Quantity < 1 {
		log.Println("cereja")
		if (item.BSServiceID != nil && *item.BSServiceID > 0) || (item.BSProductID != nil && *item.BSProductID > 0) {
			log.Println("banana")
			var orderItem entity.Item
			orderItem, err = services.ItemService.GetOrderItemByServiceProduct(orderDB.ID, item.BSServiceID, item.BSProductID)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}

			if orderItem.ID != 0 {
				err = services.ItemService.DeleteOrderItem(orderID, orderItem.ID)
				if err != nil {
					config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
					return
				}
			}
		} else {
			config.ResponsePerErr(w, err, config.INVALIDREQUEST)
			return
		}
	} else {

		if item.BSServiceID != nil && *item.BSServiceID > 0 {
			item.BSProductID = nil
			var service entity.Service
			service, err = services.BS.GetService(bsID, *item.BSServiceID)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
			////log.Printf("bs_service: %+v, bs_item: %+v", bs_service, bs_item)
			if service.ID == 0 {
				config.ResponsePerErr(w, err, config.INVALIDREQUEST)
				return
			}
			item.Name = service.Name
			item.Value = service.Value
			item.BSServiceID = &service.ID
		} else if item.BSProductID != nil && *item.BSProductID > 0 {
			item.BSServiceID = nil
			var product entity.ProductSimple
			product, err = services.ProductService.GetProductSimple(bsID, *item.BSProductID)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
			////log.Printf("bs_service: %+v, bs_item: %+v", bs_service, bs_item)
			if product.ID == 0 {
				config.ResponsePerErr(w, err, config.INVALIDREQUEST)
				return
			}
			item.Name = product.Description
			item.Value = product.SalesPrice //((bs_product.CostPrice * bs_product.ProfitMargin) / 10000) + bs_product.CostPrice
			item.BSProductID = &product.ID
		} else {
			config.ResponsePerErr(w, errors.New("produto ou serviço não enviado"), config.INVALIDREQUEST)
			return
		}

		qtt := decimal.NewFromInt(item.Quantity)
		item.TotalValue = item.Value.Mul(qtt).Sub(item.Discount)

		if item.TotalValue.IsNegative() {
			config.ResponsePerErr(w, err, config.INVALIDREQUEST)
			return
		}

		var orderItem entity.Item
		orderItem, err = services.ItemService.GetOrderItemByServiceProduct(orderDB.ID, item.BSServiceID, item.BSProductID)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		item.BSOrderID = orderID

		if orderItem.ID != 0 {
			item.ID = orderItem.ID
			err = services.ItemService.UpdateOrderItem(item)

			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		} else {
			item.ID, err = services.ItemService.CreateOrderItem(item)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}
	}

	orderItems, err := services.ItemService.GetOrdersItems(orderDB.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var totalValueTest decimal.Decimal

	orderDB.Value = decimal.Zero
	for _, v := range orderItems {
		qtt := decimal.NewFromInt(v.Quantity)
		orderDB.Value = orderDB.Value.Add(v.Value.Mul(qtt))
		orderDB.Discount = orderDB.Discount.Add(v.Discount)
		orderDB.TotalValue = orderDB.TotalValue.Add(v.TotalValue)
		totalValueTest = totalValueTest.Add(v.Value.Mul(qtt).Sub(v.Discount))
	}

	if totalValueTest != orderDB.TotalValue {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.OrderService.UpdateOrderValues(orderDB)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	respItem, err := services.ItemService.GetOrderItemFull(orderDB.ID, item.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(respItem, http.StatusOK, w)
}

func OrdersUpdateItem(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.ORDER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	var item entity.Item
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var orderDB entity.Order

	orderDB, err = services.OrderService.GetOrderIDByID(bsID, orderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if orderDB.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	if orderDB.StatusID != constants.OrderStatusOpen {
		config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
		return
	}

	if item.Quantity < 1 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	item.BSOrderID = orderID

	orderItemDB, err := services.ItemService.GetOrderItemIDByID(orderDB.ID, item.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if orderItemDB.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	var service entity.Service

	service, err = services.BS.GetService(bsID, *orderItemDB.BSServiceID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	// //log.Println(bs_service)
	// //log.Printf("%+v", bs_item)
	if service.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	qtt := decimal.NewFromInt(item.Quantity)
	item.TotalValue = service.Value.Mul(qtt).Sub(item.Discount)

	if item.TotalValue.IsNegative() {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	item.BSBMID = orderItemDB.BSBMID
	item.Name = service.Name
	item.Value = service.Value

	err = services.ItemService.UpdateOrderItem(item)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	orderItems, err := services.ItemService.GetOrdersItems(orderDB.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var totalValueTest decimal.Decimal

	orderDB.Value = decimal.Zero
	for _, v := range orderItems {
		qtt := decimal.NewFromInt(v.Quantity)
		orderDB.Value = orderDB.Value.Add(v.Value.Mul(qtt))
		orderDB.Discount = orderDB.Discount.Add(v.Discount)
		orderDB.TotalValue = orderDB.TotalValue.Add(v.TotalValue)
		totalValueTest = totalValueTest.Add(v.Value.Mul(qtt).Sub(v.Discount))
	}

	if totalValueTest != orderDB.TotalValue {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	//log.Println(orderDB)
	err = services.OrderService.UpdateOrderValues(orderDB)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	respItem, err := services.ItemService.GetOrderItemFull(orderDB.ID, item.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(respItem, http.StatusOK, w)
}

func OrdersDeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.ORDER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	itemID, err := strconv.Atoi(chi.URLParam(r, "item_id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var order entity.Order

	order, err = services.OrderService.GetOrderIDByID(bsID, orderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if order.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	if order.StatusID != constants.OrderStatusOpen {
		config.ResponsePerErr(w, err, config.CANNOTBECHANGED)
		return
	}

	err = services.ItemService.DeleteOrderItem(orderID, itemID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	orderItems, err := services.ItemService.GetOrdersItems(orderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var totalValueTest decimal.Decimal

	for _, v := range orderItems {
		order.Value = order.Value.Add(v.Value)
		order.Discount = order.Discount.Add(v.Discount)
		order.TotalValue = order.TotalValue.Add(v.TotalValue)
		qtt := decimal.NewFromInt(v.Quantity)
		totalValueTest = totalValueTest.Add(v.Value.Mul(qtt).Sub(v.Discount))
	}

	if totalValueTest != order.TotalValue {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.OrderService.UpdateOrderValues(order)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
