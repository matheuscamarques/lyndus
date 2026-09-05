package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSItemRepositoryInterface interface {
	GetOrderItemIDByID(orderID, orderItemID int) (item entity.Item, err error)
	GetOrderItemByService(bsOrderID, bsServiceID int) (item entity.Item, err error)
	GetOrderItemFull(bsOrderID, orderItemID int) (item entity.Item, err error)
	UpdateOrderItem(item entity.Item) error
	CreateOrderItem(item entity.Item) (id int, err error)
	GetOrdersItems(bsOrderID int) (items []entity.Item, err error)
	DeleteOrderItem(bsOrderId, orderItemID int) error
	GetOrderItemByServiceProduct(bsOrderID int, bsServiceID, bsProductID *int) (item entity.Item, err error)
}
