package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/utils"
)

const (
	PROFILE   = 1
	BANNER    = 2
	CATEGORY  = 3
	SERVICE   = 4
	BM        = 5
	PERSON    = 6
	ORDER     = 7
	SCHEDULE  = 8
	USER      = 9
	PRODUCT   = 10
	SUPPLIER  = 11
	BANK      = 12
	CHECKOUT  = 13
	FINANCIAL = 14

	NONE = 1
	VIEW = 2
	ALL  = 3

	Insumo  = 1
	Revenda = 2

	ProductStatusActive   = 1
	ProductStatusInactive = 2

	PaymentStatusOpen           = 1
	PaymentStatusWaitingPayment = 2
	PaymentStatusPaid           = 3
	PaymentStatusCanceled       = 4
	PaymentStatusDenied         = 5
	PaymentStatusExpired        = 6

	LyndusPayment = 1
	CreditCard    = 2
	DebitCard     = 3
	CashPayment   = 4
	FreePayment   = 5

	Scheduled   = 1
	Canceled    = 2
	Attendance  = 3
	WaitPayment = 4
	Finished    = 5
	Absent      = 6

	Site      = 1
	Facebook  = 2
	Instagram = 3
)

var ProductClasses = []entity.ProductClass{
	{ID: Insumo, Name: "Insumo"},
	{ID: Revenda, Name: "Revenda"},
}

var Permissions = []utils.Permission{
	{PROFILE, "Perfil"},
	{BANNER, "Banner"},
	{CATEGORY, "Categoria"},
	{SERVICE, "Serviço"},
	{BM, "BM"},
	{PERSON, "Cliente"},
	{ORDER, "Atendimento"},
	{SCHEDULE, "Agenda"},
	{USER, "Usuário"},
	{PRODUCT, "Produto"},
	{SUPPLIER, "Fornecedor"},
	{BANK, "Dados Bancáros"},
	{CHECKOUT, "Checkout"},
	{FINANCIAL, "Financeiro"},
}

type social struct {
	ID   int    `json:"id"`
	Desc string `json:"desc"`
}

var Socials = []social{
	{Site, "Site"},
	{Facebook, "Facebook"},
	{Instagram, "Instagram"},
}