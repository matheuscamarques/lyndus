package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSOrderRepositoryInterface interface {
	CreateOrder(order entity.Order) (id int, err error)
	GetOrders(bsID int) (orders []entity.Order, err error)
	GetOrderByID(bsID, orderID int) (order entity.Order, err error)
	GetOrderIDByID(bsID, orderID int) (order entity.Order, err error)
	GetOrdersItemsFull(orderID int) (items []entity.Item, err error)
	UpdateOrderValues(order entity.Order) error
	UpdateOrderStatus(bsID, orderID, orderStatusID int) error
}
