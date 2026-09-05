package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BSPaymentService struct {
	repo contracts.BSPaymentRepositoryInterface
}

var PaymentService contracts.BSPaymentServiceInterface

func NewBSPaymentService() *BSPaymentService {
	return &BSPaymentService{
		repo: repository.NewBSPaymentRepository(),
	}
}

func (bps BSPaymentService) GetPaymentTypes() (paymentTypes []entity.PaymentType, err error) {
	return bps.repo.GetPaymentTypes()
}
func (bps BSPaymentService) CreateOrderPayment(charge entity.PaymentCharge) (id int, err error) {
	return bps.repo.CreateOrderPayment(charge)
}
func (bps BSPaymentService) CreateOrderReceivable(receivable entity.BalanceReceivable) (id int, err error) {
	return bps.repo.CreateOrderReceivable(receivable)
}
func (bps BSPaymentService) UpdateOrderReceivableStatus(bsID, id, status int) error {
	return bps.repo.UpdateOrderReceivableStatus(bsID, id, status)
}
func (bsp BSPaymentService) ListOrderReceivable(bsID, orderID int) (receivables []entity.BalanceReceivable, err error) {
	return bsp.repo.ListOrderReceivable(bsID, orderID)
}