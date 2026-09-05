package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BSOrderService struct {
	repo contracts.BSOrderRepositoryInterface
}

var OrderService contracts.BSOrderServiceInterface

func NewBSOrderService() *BSOrderService {
	return &BSOrderService{
		repo: repository.NewBSOrderRepository(),
	}
}

func (bos BSOrderService) CreateOrder(order entity.Order) (id int, err error) {
	return bos.repo.CreateOrder(order)
}
func (bos BSOrderService) GetOrders(bsID int) (orders []entity.Order, err error) {
	return bos.repo.GetOrders(bsID)
}
func (bos BSOrderService) GetOrderByID(bsID, orderID int) (order entity.Order, err error) {
	return bos.repo.GetOrderByID(bsID, orderID)
}
func (bos BSOrderService) GetOrderIDByID(bsID, orderID int) (order entity.Order, err error) {
	return bos.repo.GetOrderIDByID(bsID, orderID)
}
func (bos BSOrderService) GetOrdersItemsFull(orderID int) (items []entity.Item, err error) {
	return bos.repo.GetOrdersItemsFull(orderID)
}
func (bos BSOrderService) UpdateOrderValues(order entity.Order) error {
	return bos.repo.UpdateOrderValues(order)
}
func (bos BSOrderService) UpdateOrderStatus(bsID, orderID, orderStatusID int) error {
	return bos.repo.UpdateOrderStatus(bsID, orderID, orderStatusID)
}
