package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type BSItemRepository struct {
	contracts.BSItemRepositoryInterface
	conn *sqlx.DB
}

func NewBSItemRepository() BSItemRepository {
	return BSItemRepository{
		conn: postgres.DB,
	}
}

func (i BSItemRepository) GetOrderItemIDByID(orderID, orderItemID int) (item entity.Item, err error) {
	command := `SELECT boi.id,
                  	   boi.bs_order_id,
       				   boi.bs_service_id,
					   boi.bs_product_id
                FROM bs_order_item boi  
                WHERE boi.bs_order_id=$1 AND boi.id=$2`
	err = i.conn.Get(&item, command, orderID, orderItemID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return item, err
}

func (i BSItemRepository) GetOrderItemByService(bsOrderID, bsServiceID int) (item entity.Item, err error) {
	command := `SELECT id,
       				   bs_order_id,
       				   bs_service_id,
       				   bs_product_id,
       				   bs_bm_id,
       				   name,
       				   value,
       				   discount,
       				   total_value,
       				   quantity 
            FROM bs_order_item WHERE bs_order_id=$1 AND bs_service_id=$2`
	err = i.conn.Get(&item, command, bsOrderID, bsServiceID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return item, err
}

func (i BSItemRepository) GetOrderItemByServiceProduct(bsOrderID int, bsServiceID, bsProductID *int) (item entity.Item, err error) {
	command := `SELECT id,
       				   bs_order_id,
       				   bs_service_id,
       				   bs_product_id,
       				   bs_bm_id,
       				   name,
       				   value,
       				   discount,
       				   total_value,
       				   quantity 
            FROM bs_order_item WHERE bs_order_id=$1 AND (bs_service_id=$2 OR bs_product_id=$3)`
	err = i.conn.Get(&item, command, bsOrderID, bsServiceID, bsProductID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return item, err
}

func (i BSItemRepository) GetOrderItemFull(bsOrderID, orderItemID int) (item entity.Item, err error) {
	command := `SELECT boi.id,
					   boi.bs_service_id,
       				   boi.bs_product_id,
					   boi.bs_bm_id,
					   bs_bm.name as bs_bm_name, 
					   boi.name,
					   boi.value,
					   boi.discount,
					   boi.total_value,
					   boi.quantity 
                FROM bs_order_item boi
                LEFT JOIN bs_bm bs_bm ON boi.bs_bm_id=bs_bm.id
                WHERE boi.bs_order_id = $1 AND boi.id = $2`
	err = i.conn.Get(&item, command, bsOrderID, orderItemID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return item, err
}

func (i BSItemRepository) UpdateOrderItem(item entity.Item) error {
	command := `UPDATE bs_order_item SET name=$3,
                         				 value=$4,
                         				 discount=$5,
                         				 total_value=$6,
                         				 quantity=$7
        		WHERE bs_order_id=$1 AND id=$2`
	stmt, err := i.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.BSOrderID, item.ID, item.Name, item.Value, item.Discount, item.TotalValue, item.Quantity)
	return err
}

func (i BSItemRepository) CreateOrderItem(item entity.Item) (id int, err error) {
	command := `INSERT INTO bs_order_item(bs_order_id, 
                          				  bs_service_id,
                          				  bs_product_id,
                          				  bs_bm_id,
                          				  name,
                          				  value,
                          				  discount,
                          				  quantity,
                          				  total_value)
                          		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
                          		RETURNING ID`
	stmt, err := i.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(item.BSOrderID,
		item.BSServiceID,
		item.BSProductID,
		item.BSBMID,
		item.Name,
		item.Value,
		item.Discount,
		item.Quantity,
		item.TotalValue).Scan(&id)
	return id, err
}

func (i BSItemRepository) GetOrdersItems(bsOrderID int) (items []entity.Item, err error) {
	command := `SELECT id,
					   bs_order_id,
					   bs_service_id,
       				   bs_product_id,
					   bs_bm_id,
					   name,
					   value,
					   discount,
					   total_value,
					   quantity 
				FROM bs_order_item
				WHERE bs_order_id=$1 `
	err = i.conn.Select(&items, command, bsOrderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return items, err
}

func (i BSItemRepository) DeleteOrderItem(bsOrderId, orderItemID int) error {
	command := `DELETE FROM bs_order_item WHERE bs_order_id = $1 AND id = $2`
	stmt, err := i.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsOrderId, orderItemID)
	return err
}
