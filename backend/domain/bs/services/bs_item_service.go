package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BSItemService struct {
	repo contracts.BSItemRepositoryInterface
}

var ItemService contracts.BSItemServiceInterface

func NewBSItemService() *BSItemService {
	return &BSItemService{
		repo: repository.NewBSItemRepository(),
	}
}

func (bis BSItemService) GetOrderItemIDByID(orderID, orderItemID int) (item entity.Item, err error) {
	return bis.repo.GetOrderItemIDByID(orderID, orderItemID)
}
func (bis BSItemService) GetOrderItemByService(bsOrderID, bsServiceID int) (item entity.Item, err error) {
	return bis.repo.GetOrderItemByService(bsOrderID, bsServiceID)
}
func (bis BSItemService) GetOrderItemFull(bsOrderID, orderItemID int) (item entity.Item, err error) {
	return bis.repo.GetOrderItemFull(bsOrderID, orderItemID)
}
func (bis BSItemService) UpdateOrderItem(item entity.Item) error {
	return bis.repo.UpdateOrderItem(item)
}
func (bis BSItemService) CreateOrderItem(item entity.Item) (id int, err error) {
	return bis.repo.CreateOrderItem(item)
}
func (bis BSItemService) GetOrdersItems(bsOrderID int) (items []entity.Item, err error) {
	return bis.repo.GetOrdersItems(bsOrderID)
}
func (bis BSItemService) DeleteOrderItem(bsOrderId, orderItemID int) error {
	return bis.repo.DeleteOrderItem(bsOrderId, orderItemID)
}
func (bis BSItemService) GetOrderItemByServiceProduct(bsOrderID int, bsServiceID, bsProductID *int) (item entity.Item, err error) {
	return bis.repo.GetOrderItemByServiceProduct(bsOrderID, bsServiceID, bsProductID)
}