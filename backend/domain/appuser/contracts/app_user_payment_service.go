package contracts

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
)

type AppUserPaymentInterface interface {
	GetOpenPayments(appUserID, paymentStatusID int) (payments []entity.Payment, err error)
	GetOrderByID(orderID int) (order entity.Order, err error)
	GetOrderItems(orderID int) (items []entity.Item, err error)
	GetOrderPayments(bsOrderID int) (payments []entity.OrderPayments, err error)
	GetBsName(bsID int) (name string, err error)
	GetPayment(paymentUUID string) (confirmPayment entity.ConfirmPayment, err error)

	UpdatePayment(id, bsID, paymentStatusID int) error
	InsertStatement(statement entity.AppUserStatement) (int, error)
	UpdateBsOrderStatus(bsOrderID, bsID, orderStatusID int) error
	GetOrderID(appUserID, bsOrderPaymentID int) (bsOrderID int, err error)
	UpdateScheduleStatus(bsOrderID, bsID, scheduleStatusID int) error

	SaveCredCard(appUserID int, paymentTypeCode, token, maskedCardNumber string) (err error)
}
