package config

const (
	CLIENTPROFILE  = 1
	CLIENTCATEGORY = 2
	CLIENTEMPLOYEE = 3
	CLIENTBENEFIT  = 4
	CLIENTUSER     = 5
	CLIENTPOLL     = 6
)

const (
	BSPROFILE  = 1
	BSBANNER   = 2
	BSCATEGORY = 3
	BSSERVICE  = 4
	BSBM       = 5
	BSPERSON   = 6
	BSORDER    = 7
	BSSCHEDULE = 8
	BSUSER     = 9
)

var BSPERMISSIONS = map[int]string{
	BSPROFILE:  "Perfil",
	BSBANNER:   "Banner",
	BSCATEGORY: "Categoria",
	BSSERVICE:  "Serviço",
	BSBM:       "BM",
	BSPERSON:   "Cliente",
	BSORDER:    "Pedido",
	BSSCHEDULE: "Agenda",
	BSUSER:     "Usuário",
}

var CLIENTPERMISSIONS = map[int]string{
	CLIENTPROFILE:  "Perfil",
	CLIENTCATEGORY: "Categoria",
	CLIENTEMPLOYEE: "Funcionário",
	CLIENTBENEFIT:  "Benefício",
	CLIENTUSER:     "Usuário",
	CLIENTPOLL:		"Enquete",
}

var LEVEL = map[int]string{
	1: "Nenhuma",
	2: "Visualizar",
	3: "Todas",
}
