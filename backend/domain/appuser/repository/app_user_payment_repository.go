package repository

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type AppUserPaymentRepository struct {
	contracts.AppUserPaymentRepositoryInterface
	conn *sqlx.DB
}

func NewAppUserPaymentRepository() *AppUserPaymentRepository {
	return &AppUserPaymentRepository{
		conn: postgres.DB,
	}
}

func (aus AppUserPaymentRepository) SaveCredCard(appUserID int, paymentTypeCode, token, maskedCardNumber string) (err error) {
	command := `INSERT INTO app_user_cards(app_user_id,
                           				   payment_type_code,
                           				   token,
                           				   masked_card_number)
						VALUES($1,$2,$3,$4) RETURNING id`
	stmt, err := aus.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(appUserID, paymentTypeCode, token, maskedCardNumber)
	return err

}

func (aus AppUserPaymentRepository) GetOpenPayments(appUserID, paymentStatusID int) (payments []entity.Payment, err error) {

	command := `select bp.app_user_id as id,
       				   bbr.payment_uuid,
       				   bbr.value,
       				   bop.bs_order_id
					FROM bs_balance_receivable bbr
						INNER JOIN  bs_person bp ON bbr.bs_person_id = bp.id
						INNER JOIN bs_order_payment bop ON bbr.bs_order_payment_id = bop.id
						WHERE bp.app_user_id = $1 
						  AND bbr.payment_status_id = $2`
	err = aus.conn.Select(&payments, command, appUserID, paymentStatusID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return payments, err
}

func (aus AppUserPaymentRepository) GetOrderID(appUserID, bsOrderPaymentID int) (bsOrderID int, err error) {
	command := `SELECT bop.bs_order_id
					FROM bs_balance_receivable bbr
					    INNER JOIN  bs_person bp ON bbr.bs_person_id = bp.id
						INNER JOIN bs_order_payment bop ON bbr.bs_order_payment_id = bop.id
						WHERE bp.app_user_id = $1 
						  AND bop.id = $2`
	err = aus.conn.Get(&bsOrderID, command, appUserID, bsOrderPaymentID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bsOrderID, err
}

func (aus AppUserPaymentRepository) GetOrderByID(orderID int) (order entity.Order, err error) {
	command := `SELECT bso.bs_id,
       				   bso.value,
					   bso.discount,
					   bso.total_value,
					   bso.created_at,
					   ost."desc_pt" as status
                FROM bs_order bso
                INNER JOIN order_status ost ON bso.order_status_id = ost.id   
                WHERE bso.id=$1`
	err = aus.conn.Get(&order, command, orderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return order, err
}

func (aus AppUserPaymentRepository) GetOrderItems(orderID int) (items []entity.Item, err error) {
	command := `SELECT bs_bm.name as bs_bm_name, 
					   boi.name,
					   boi.value,
					   boi.discount,
					   boi.total_value,
					   boi.quantity 
                FROM bs_order_item boi
                LEFT JOIN bs_bm bs_bm ON boi.bs_bm_id=bs_bm.id
                WHERE boi.bs_order_id = $1 `
	err = aus.conn.Select(&items, command, orderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return items, err
}

func (aus AppUserPaymentRepository) GetOrderPayments(bsOrderID int) (payments []entity.OrderPayments, err error) {
	command := `SELECT pt.nome as name,
       				   bop.value
                FROM bs_order_payment bop
         		LEFT JOIN payment_type pt ON pt.id = bop.payment_type_id
                	WHERE bop.bs_order_id = $1`
	err = aus.conn.Select(&payments, command, bsOrderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return payments, err
}

func (aus AppUserPaymentRepository) GetBsName(bsID int) (name string, err error) {
	command := `SELECT com.company_name 
					FROM bs
        			INNER JOIN company com on com.id = bs.company_id
            			WHERE bs.id=$1`
	err = aus.conn.Get(&name, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return name, err
}

func (aus AppUserPaymentRepository) GetPayment(paymentUUID string) (confirmPayment entity.ConfirmPayment, err error) {
	command := `SELECT id,
       				   bs_id,
       				   bs_order_payment_id,
       				   value,
       				   payment_uuid
					FROM bs_balance_receivable bbr
            		WHERE bbr.payment_uuid=$1`
	err = aus.conn.Get(&confirmPayment, command, paymentUUID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return confirmPayment, err
}

func (aus AppUserPaymentRepository) UpdatePayment(id, bsID, paymentStatusID int) error {
	command := `UPDATE bs_balance_receivable SET payment_status_id=$3 WHERE id=$1 AND bs_id=$2`
	stmt, err := aus.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, bsID, paymentStatusID)
	return err
}

func (aus AppUserPaymentRepository) UpdateBsOrderStatus(bsOrderID, bsID, orderStatusID int) error {
	command := `UPDATE bs_order SET order_status_id=$3 WHERE id=$1 AND bs_id=$2`
	stmt, err := aus.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsOrderID, bsID, orderStatusID)
	return err
}

func (aus AppUserPaymentRepository) UpdateScheduleStatus(bsOrderID, bsID, scheduleStatusID int) error {
	command := `UPDATE bs_schedule_time SET schedule_status_id=$3 WHERE bs_order_id=$1 AND bs_id=$2`
	stmt, err := aus.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsOrderID, bsID, scheduleStatusID)
	return err
}

func (aus AppUserPaymentRepository) InsertStatement(statement entity.AppUserStatement) (int, error) {
	command := `INSERT INTO app_user_statement(app_user_balance_id,
                               				   app_user_id,                               				   
                               				   value,
                               				   statements_id,
                               				   "desc",
                               				   date_time)
						VALUES($1,$2,$3,$4,$5,now()) RETURNING id`
	stmt, err := aus.conn.Prepare(command)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(statement.AppUserBalanceID,
		statement.AppUserID,
		statement.Value,
		statement.StatementID,
		statement.Desc).Scan(&id)
	return id, err
}
