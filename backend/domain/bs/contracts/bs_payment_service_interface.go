package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSPaymentServiceInterface interface {
	GetPaymentTypes() (paymentTypes []entity.PaymentType, err error)
	CreateOrderPayment(charge entity.PaymentCharge) (id int, err error)
	CreateOrderReceivable(receivable entity.BalanceReceivable) (id int, err error)
	ListOrderReceivable(bsID, orderID int) (receivables []entity.BalanceReceivable, err error)
	UpdateOrderReceivableStatus(bsID, id, status int) error
}
