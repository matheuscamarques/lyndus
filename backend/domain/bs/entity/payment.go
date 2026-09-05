package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type PaymentType struct {
	ID   int    `json:"id" db:"id"`
	Desc string `json:"desc" db:"desc"`
}

type PaymentCharge struct {
	ID            int             `json:"id" db:"id"`
	BSOrderID     int             `json:"bsOrderID" db:"bs_order_id"`
	PaymentTypeID int             `json:"paymentTypeID" db:"payment_type_id"`
	Value         decimal.Decimal `json:"value" db:"value"`
}

type BalanceReceivable struct {
	ID               int             `json:"id" db:"id"`
	BsID             int             `json:"bsID" db:"bs_id"`
	BSPersonID       int             `json:"BSPersonID_id" db:"bs_person_id"`
	BSOrderPaymentID int             `json:"BSOrderPaymentID" db:"bs_order_payment_id"`
	PaymentStatusID  int             `json:"paymentStatusID" db:"payment_status_id"`
	PaymentUUID      uuid.UUID       `json:"paymentUUID" db:"payment_uuid"`
	Value            decimal.Decimal `json:"value" db:"value"`

	BSWithdrawalCashID int       `json:"BSWithdrawalCashID" db:"bs_withdrawal_cash_id"`
	CreatedAT          time.Time `json:"createdAT" db:"created_at"`
}
