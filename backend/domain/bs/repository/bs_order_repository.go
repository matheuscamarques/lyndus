package repository

import (
	"bitbucket.org/lyndus/backend/domain/constants"
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type BSOrderRepository struct {
	contracts.BSOrderRepositoryInterface
	conn *sqlx.DB
}

func NewBSOrderRepository() BSOrderRepository {
	return BSOrderRepository{
		conn: postgres.DB,
	}
}

func (o BSOrderRepository) CreateOrder(order entity.Order) (id int, err error) {
	command := `INSERT INTO bs_order(bs_id,
                     				 order_status_id,
                     				 bs_person_id,
                     				 value,
                     				 discount,
                     				 total_value,
                     				 obs)
        			VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING ID`
	stmt, err := o.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(order.BsID,
		order.StatusID,
		order.BSPersonID,
		order.Value,
		order.Discount,
		order.TotalValue,
		order.Obs).Scan(&id)
	return id, err
}

func (o BSOrderRepository) GetOrders(bsID int) (orders []entity.Order, err error) {
	command := `SELECT bso.id,
					   bso.bs_person_id,
					   per.name as bs_person_name,
					   ors."desc_pt" as status,
					   bso.value,
					   bso.discount,
					   bso.total_value,
					   bso.created_at
				FROM bs_order bso
				INNER JOIN bs_person per ON per.id=bso.bs_person_id
				INNER JOIN order_status ors ON bso.order_status_id=ors.id
				WHERE bso.bs_id=$1 ORDER BY created_at DESC`
	err = o.conn.Select(&orders, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return orders, err
}

func (o BSOrderRepository) GetOrderByID(bsID, orderID int) (order entity.Order, err error) {
	command := `SELECT bso.id,
                  	   bso.bs_id,
                  	   bso.bs_person_id,
                  	   per.name as bs_person_name,
                  	   bso.order_status_id,
					   bso.value,
					   bso.discount,
					   bso.total_value,
					   bso.obs,
					   bso.created_at
                FROM bs_order bso
                INNER JOIN bs_person per ON per.id = bso.bs_person_id   
                WHERE bso.bs_id=$1 AND bso.id=$2`
	err = o.conn.Get(&order, command, bsID, orderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return order, err
}

func (o BSOrderRepository) GetOrderIDByID(bsID, orderID int) (order entity.Order, err error) {
	command := `SELECT bso.id,
                  	   bso.bs_id,
       				   bso.order_status_id
                FROM bs_order bso  
                WHERE bso.bs_id=$1 AND bso.id=$2`
	err = o.conn.Get(&order, command, bsID, orderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return order, err
}

func (o BSOrderRepository) GetOrdersItemsFull(orderID int) (items []entity.Item, err error) {
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
                WHERE boi.bs_order_id = $1
				ORDER BY boi.name`
	err = o.conn.Select(&items, command, orderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return items, err
}

func (o BSOrderRepository) UpdateOrderValues(order entity.Order) error {
	command := `UPDATE bs_order SET value=$3,
                    				discount=$4,
                    				total_value=$5
					WHERE bs_id=$1 AND id=$2 AND order_status_id <> $6`
	stmt, err := o.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(order.BsID, order.ID, order.Value, order.Discount, order.TotalValue, constants.OrderStatusCanceled)
	return err
}

func (o BSOrderRepository) UpdateOrderStatus(bsID, orderID, orderStatusID int) error {
	command := `UPDATE bs_order SET order_status_id=$3 WHERE bs_id=$1 AND id=$2`
	stmt, err := o.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, orderID, orderStatusID)
	return err
}
