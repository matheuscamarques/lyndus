package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
)

type Item struct {
	//BSServiceID *int          `json:"serviceID" db:"bs_service_id"`
	//BSProductID *int          `json:"productID" db:"bs_product_id"`
	BSBMName   string          `json:"bmName" db:"bs_bm_name"`
	Name       string          `json:"name" db:"name"`
	Value      decimal.Decimal `json:"value" db:"value"`
	Discount   decimal.Decimal `json:"discount" db:"discount"`
	TotalValue decimal.Decimal `json:"totalValue" db:"total_value"`
	Quantity   int64           `json:"quantity" db:"quantity"`
}

type OrderPayments struct {
	Name  string          `json:"name" db:"name"`
	Value decimal.Decimal `json:"value" db:"value"`
}

type Order struct {
	BsID       int             `json:"-" db:"bs_id"`
	Value      decimal.Decimal `json:"value" db:"value"`
	Discount   decimal.Decimal `json:"discount" db:"discount"`
	TotalValue decimal.Decimal `json:"totalValue" db:"total_value"`
	//Obs        string        `json:"obs" db:"obs"`
	Items []Item `json:"items" db:"items"`

	StatusID  int            `json:"statusID,omitempty" db:"order_status_id"`
	Status    string         `json:"status" db:"status"`
	CreatedAT types.DateTime `json:"createdAT" db:"created_at"`

	Payments []OrderPayments `json:"orderPayments" db:"order_payments"`
}
type Payment struct {
	ID          int             `json:"-" db:"id"`
	BSName      string          `json:"bsName" json:"bs_name"`
	BsID        *int            `json:"bsID"`
	PaymentUUID string          `json:"paymentUUID" db:"payment_uuid"`
	Value       decimal.Decimal `json:"value" db:"value"`
	Order       Order           `json:"order" db:"order"`
	//Status      string        `json:"status" db:"status"`
	BSOrderID int `json:"-" db:"bs_order_id"`
}

type BsOrderPayment struct {
	ID        int `json:"-" db:"id"`
	BsOrderID int `db:"bs_order_id"`

	BSName      string          `json:"bsName" json:"bs_name"`
	PaymentUUID string          `json:"paymentUUID" db:"payment_uuid"`
	Value       decimal.Decimal `json:"value" db:"value"`
	Order       Order           `json:"order" db:"order"`
	//Status      string        `json:"status" db:"status"`
	BSOrderID int `json:"-" db:"bs_order_id"`
}

type ConfirmPayment struct {
	ID               int             `json:"-" db:"id"`
	BsID             int             `json:"-" db:"bs_id"`
	BsOrderPaymentID int             `json:"-" db:"bs_order_payment_id"`
	PaymentUUID      string          `json:"paymentUUID" db:"payment_uuid"`
	Value            decimal.Decimal `json:"value" db:"value"`
	Password         string          `json:"password" db:"password"`
}

type AppUserStatement struct {
	ID               int             `json:"id" db:"id"`
	AppUserID        int             `json:"appUserID" db:"app_user_id"`
	AppUserBalanceID int             `json:"appUserBalanceID" db:"app_user_balance_id"`
	StatementID      int             `json:"statementID" db:"statement_id"`
	DateTime         types.DateTime  `json:"dateTime" db:"date_time"`
	Value            decimal.Decimal `json:"value" db:"value"`
	Desc             string          `json:"desc" db:"desc"`
}
