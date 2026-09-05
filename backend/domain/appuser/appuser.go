package appuser

const (
	BMStatusActive   = 1
	BMStatusInactive = 2

	OrderStatusOpen           = 1
	OrderStatusWaitingPayment = 2
	OrderStatusPaid           = 3
	OrderStatusCanceled       = 4
	OrderStatusExpired        = 5



	LyndusPayment = 1
	CreditCard    = 2
	DebitCard     = 3
	CashPayment   = 4
	FreePayment   = 5

	StatementBenefit = 1
	StatementPayment = 2
	StatementBonus  = 3
	StatementCredit  = 4
)

var Statements = map[int]string{
	StatementBenefit: "Benefício",
	StatementPayment: "Pagamento",
	StatementBonus:  "Bônus",
	StatementCredit:  "Crédito",
}
