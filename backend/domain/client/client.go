package client

import "bitbucket.org/lyndus/backend/infra/utils"

const (
	POLLCREATED uint8 = iota + 1
	POLLINPROGRESS
	POLLCLOSED
	POLLCANCELLED
)

const (
	PROFILE  = 1
	CATEGORY = 2
	EMPLOYEE = 3
	BENEFIT  = 4
	USER     = 5
	POLL     = 6
	NONE     = 1
	VIEW     = 2
	ALL      = 3

	EMPLOYEESTATUSINACTIVE = 2
	EMPLOYEESTATUSACTIVE   = 1

	BENEFITSTATUSOPEN                 = 1
	BENEFITSTATUSWAITINGPAYMENT       = 2
	BENEFITSTATUSPAID                 = 3
	BENEFITSTATUSCANCELED             = 4
	BENEFITSTATUSEXPIRED              = 5
	BENEFITSTATUSPAYMENTCONFIRMED     = 6
	BENEFITSTATUSREQUESTBILLINGTICKET = 7
)

var BENEFITSTATUSPT = map[int]string{
	BENEFITSTATUSOPEN:                 "Aberto",
	BENEFITSTATUSWAITINGPAYMENT:       "Aguardando Pagamento",
	BENEFITSTATUSPAID:                 "Pago",
	BENEFITSTATUSCANCELED:             "Cancelado",
	BENEFITSTATUSEXPIRED:              "Expirado",
	BENEFITSTATUSPAYMENTCONFIRMED:     "Pagamento Confirmado",
	BENEFITSTATUSREQUESTBILLINGTICKET: "Aguardando Boleto",
}

var Permissions = []utils.Permission{
	{PROFILE, "Perfil"},
	{CATEGORY, "Categoria"},
	{EMPLOYEE, "Funcionário"},
	{BENEFIT, "Benefício"},
	{USER, "Usuário"},
	{POLL, "Enquete"},
}
