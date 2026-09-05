package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
)

type AppUserPaymentService struct {
	repo contracts.AppUserPaymentRepositoryInterface
}

var PaymentService contracts.AppUserPaymentInterface

func NewAppUserPaymentService() *AppUserPaymentService {
	return &AppUserPaymentService{
		repo: repository.NewAppUserPaymentRepository(),
	}
}

func (ass AppUserPaymentService) GetOpenPayments(appUserID, paymentStatusID int) (payments []entity.Payment, err error) {
	return ass.repo.GetOpenPayments(appUserID, paymentStatusID)
}

func (ass AppUserPaymentService) GetOrderByID(orderID int) (order entity.Order, err error) {
	return ass.repo.GetOrderByID(orderID)
}
func (ass AppUserPaymentService) GetOrderItems(orderID int) (items []entity.Item, err error) {
	return ass.repo.GetOrderItems(orderID)
}
func (ass AppUserPaymentService) GetOrderPayments(bsOrderID int) (payments []entity.OrderPayments, err error) {
	return ass.repo.GetOrderPayments(bsOrderID)
}

func (ass AppUserPaymentService) GetBsName(bsID int) (name string, err error) {
	return ass.repo.GetBsName(bsID)
}

func (ass AppUserPaymentService) GetPayment(paymentUUID string) (confirmPayment entity.ConfirmPayment, err error) {
	return ass.repo.GetPayment(paymentUUID)
}

func (ass AppUserPaymentService) UpdatePayment(id, bsID, paymentStatusID int) error {
	return ass.repo.UpdatePayment(id, bsID, paymentStatusID)
}

func (ass AppUserPaymentService) InsertStatement(statement entity.AppUserStatement) (int, error) {
	return ass.repo.InsertStatement(statement)
}

func (ass AppUserPaymentService) UpdateBsOrderStatus(bsOrderID, bsID, orderStatusID int) error {
	return ass.repo.UpdateBsOrderStatus(bsOrderID, bsID, orderStatusID)
}

func (ass AppUserPaymentService) GetOrderID(appUserID, bsOrderPaymentID int) (bsOrderID int, err error) {
	return ass.repo.GetOrderID(appUserID, bsOrderPaymentID)
}

func (ass AppUserPaymentService) UpdateScheduleStatus(bsOrderID, bsID, scheduleStatusID int) error {
	return ass.repo.UpdateScheduleStatus(bsOrderID, bsID, scheduleStatusID)
}

func (ass AppUserPaymentService) SaveCredCard(appUserID int, paymentTypeCode, token, maskedCardNumber string) (err error) {
	return ass.repo.SaveCredCard(appUserID, paymentTypeCode, token, maskedCardNumber)
}

//func (ass AppUserPaymentService) GetServicesByCategory(bsID, serviceCategoryID int) (services []entity.Service, err error){
//	return ass.repo.GetServicesByCategory(bsID, serviceCategoryID)
//}
