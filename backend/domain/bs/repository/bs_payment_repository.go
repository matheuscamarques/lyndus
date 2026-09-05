package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type BSPaymentRepository struct {
	contracts.BSPaymentRepositoryInterface
	conn *sqlx.DB
}

func NewBSPaymentRepository() BSPaymentRepository {
	return BSPaymentRepository{
		conn: postgres.DB,
	}
}

func (p BSPaymentRepository) GetPaymentTypes() (paymentTypes []entity.PaymentType, err error) {
	command := `SELECT pty.id, pty.nome as "desc"  FROM payment_type pty order by id`
	err = p.conn.Select(&paymentTypes, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return paymentTypes, err
}

func (p BSPaymentRepository) CreateOrderPayment(charge entity.PaymentCharge) (id int, err error) {
	command := `INSERT INTO bs_order_payment(bs_order_id, payment_type_id, value) VALUES($1, $2, $3) RETURNING ID`
	stmt, err := p.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(charge.BSOrderID, charge.PaymentTypeID, charge.Value).Scan(&id)
	return id, err
}

func (p BSPaymentRepository) ListOrderReceivable(bsID, orderID int) (receivables []entity.BalanceReceivable, err error) {
	command := `SELECT bbr.bs_id,
       				   bbr.bs_person_id,
       				   bbr.bs_order_payment_id,
       				   bbr.payment_status_id,
       				   bbr.payment_uuid,
       				   bbr.value
       				FROM bs_balance_receivable bbr
					INNER JOIN bs_order_payment bop on bbr.bs_order_payment_id = bop.id
					WHERE bbr.bs_id=$1 AND bop.bs_order_id=$2
					order by id`
	err = p.conn.Select(&receivables, command, bsID, orderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return receivables, err
}

func (p BSPaymentRepository) CreateOrderReceivable(receivable entity.BalanceReceivable) (id int, err error) {
	command := `INSERT INTO bs_balance_receivable(bs_id, bs_person_id, bs_order_payment_id, payment_status_id,
        			payment_uuid, value) VALUES($1, $2, $3, $4, $5, $6) RETURNING ID`
	stmt, err := p.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(receivable.BsID, receivable.BSPersonID, receivable.BSOrderPaymentID, receivable.PaymentStatusID, receivable.PaymentUUID, receivable.Value).Scan(&id)
	return id, err
}

func (o BSOrderRepository) UpdateOrderReceivableStatus(bsID, id, status int) error {
	command := `UPDATE bs_balance_receivable SET payment_status_id=$3
					WHERE bs_id=$1 AND id=$2 `
	stmt, err := o.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, id, status)
	return err
}