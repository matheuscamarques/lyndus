package config

import "net/http"

const (
	UNAUTHORIZED           = http.StatusUnauthorized
	INTERNALSERVERERROR    = http.StatusInternalServerError
	INVALIDREQUEST         = http.StatusBadRequest
	NOTFOUND               = http.StatusNotFound
	FORBIDDEN              = http.StatusForbidden
	INVALIDMAXSIZE         = 40
	EXPIRED                = 41
	CANNOTBECHANGED        = 42
	PHONEALREADYREGISTERED = 43
	CPFINVALID             = 45
	CPFALREADYREGISTERED   = 46
	CNPJINVALID            = 47
	CNPJALREADYREGISTERED  = 48
	DOCINVALID             = 49
	DOCALREADYREGISTERED   = 50
	LOGINCHECKERROR        = 51
	INVALIDLOGIN           = 52
	INVALIDPASSWORD        = 53
	DATAINVALID            = 54
	FILEINVALID            = 55
	ORDERCANCELED          = 56
	ORDEREXPIRED           = 57
	ORDERINVALIDVALUE      = 58
	SCHEDULEOVERLAY        = 59
	SCHEDULEOVERLAYOPEN    = 60

	WRONGVALUE          = 65
	BALANCEINSUFFICIENT = 66
	PAYMENTPASSWORD     = 67

	USERNOTFOUND  = 100
	CANNOTBESAVED = 101

	USERALREADYEXISTS = 200
)

var errMap = map[int]Err{
	UNAUTHORIZED: {
		ResponseErr: ResponseErr{
			Code:    UNAUTHORIZED,
			Message: "Acesso inválido",
		},
		HttpCode: UNAUTHORIZED,
	},
	INTERNALSERVERERROR: {
		ResponseErr: ResponseErr{
			Code:    INTERNALSERVERERROR,
			Message: "Erro Interno, se persistir contate o suporte.",
		},
		HttpCode: INTERNALSERVERERROR,
	},
	NOTFOUND: {
		ResponseErr: ResponseErr{
			Code:    NOTFOUND,
			Message: "Não foi possível localizar.",
		},
		HttpCode: NOTFOUND,
	},
	FORBIDDEN: {
		ResponseErr: ResponseErr{
			Code:    FORBIDDEN,
			Message: "Você não tem permissão para executar essa ação..",
		},
		HttpCode: FORBIDDEN,
	},

	INVALIDREQUEST: {
		ResponseErr: ResponseErr{
			Message: "Solicitação inválida.",
			Code:    INVALIDREQUEST,
		},
		HttpCode: INVALIDREQUEST,
	},
	INVALIDMAXSIZE: {
		ResponseErr: ResponseErr{
			Message: "Limite máximo de upload é 20MB.",
			Code:    INVALIDMAXSIZE,
		},
		HttpCode: http.StatusBadRequest,
	},
	EXPIRED: {
		ResponseErr: ResponseErr{
			Message: "Solicitação inválida ou expirada.",
			Code:    EXPIRED,
		},
		HttpCode: http.StatusBadRequest,
	},
	CPFINVALID: {
		ResponseErr: ResponseErr{
			Message: "CPF inválido.",
			Code:    CPFINVALID,
		},
		HttpCode: http.StatusBadRequest,
	},
	CNPJINVALID: {
		ResponseErr: ResponseErr{
			Message: "CNPJ inválido.",
			Code:    CNPJINVALID,
		},
		HttpCode: http.StatusBadRequest,
	},
	DATAINVALID: {
		ResponseErr: ResponseErr{
			Message: "Dados enviados inválidos.",
			Code:    DATAINVALID,
		},
		HttpCode: http.StatusBadRequest,
	},
	FILEINVALID: {
		ResponseErr: ResponseErr{
			Message: "Arquivo enviado náo é válido.",
			Code:    FILEINVALID,
		},
		HttpCode: http.StatusBadRequest,
	},

	DOCINVALID: {
		ResponseErr: ResponseErr{
			Message: "Documento inválido.",
			Code:    DOCINVALID,
		},
		HttpCode: http.StatusBadRequest,
	},
	DOCALREADYREGISTERED: {
		ResponseErr: ResponseErr{
			Message: "Documento já cadastrado.",
			Code:    DOCALREADYREGISTERED,
		},
		HttpCode: http.StatusBadRequest,
	},
	PHONEALREADYREGISTERED: {
		ResponseErr: ResponseErr{
			Message: "Telefone já cadastrado.",
			Code:    PHONEALREADYREGISTERED,
		},
		HttpCode: http.StatusBadRequest,
	},
	CPFALREADYREGISTERED: {
		ResponseErr: ResponseErr{
			Message: "CPF já cadastrado.",
			Code:    CPFALREADYREGISTERED,
		},
		HttpCode: http.StatusBadRequest,
	},
	CNPJALREADYREGISTERED: {
		ResponseErr: ResponseErr{
			Message: "CNPJ já cadastrado.",
			Code:    CNPJALREADYREGISTERED,
		},
		HttpCode: http.StatusBadRequest,
	},

	CANNOTBECHANGED: {
		ResponseErr: ResponseErr{
			Message: "Não pode ser alterado.",
			Code:    CANNOTBECHANGED,
		},
		HttpCode: http.StatusBadRequest,
	},
	CANNOTBESAVED: {
		ResponseErr: ResponseErr{
			Message: "Falha ao salvar.",
			Code:    CANNOTBESAVED,
		},
		HttpCode: http.StatusInternalServerError,
	},
	LOGINCHECKERROR: {
		ResponseErr: ResponseErr{
			Message: "Erro ao verificar login",
			Code:    LOGINCHECKERROR,
		},
		HttpCode: http.StatusInternalServerError,
	},
	INVALIDPASSWORD: {
		ResponseErr: ResponseErr{
			Message: "Erro ao verificar login",
			Code:    INVALIDPASSWORD,
		},
		HttpCode: http.StatusBadRequest,
	},

	INVALIDLOGIN: {
		ResponseErr: ResponseErr{
			Message: "Login ou senha inválidos.",
			Code:    LOGINCHECKERROR,
		},
		HttpCode: http.StatusBadRequest,
	},
	USERNOTFOUND: {
		ResponseErr: ResponseErr{
			Message: "Usuário não encontrado.",
			Code:    USERNOTFOUND,
		},
		HttpCode: http.StatusNotFound,
	},
	ORDERCANCELED: {
		ResponseErr: ResponseErr{
			Message: "Pedido cancelado.",
			Code:    ORDERCANCELED,
		},
		HttpCode: http.StatusBadRequest,
	},
	ORDEREXPIRED: {
		ResponseErr: ResponseErr{
			Message: "Pedido expirado.",
			Code:    ORDEREXPIRED,
		},
		HttpCode: http.StatusBadRequest,
	},
	ORDERINVALIDVALUE: {
		ResponseErr: ResponseErr{
			Message: "Valor informado inválido.",
			Code:    ORDERINVALIDVALUE,
		},
		HttpCode: http.StatusBadRequest,
	},
	SCHEDULEOVERLAY: {
		ResponseErr: ResponseErr{
			Message: "Agenda sobreposta a horário já agendado.",
			Code:    SCHEDULEOVERLAY,
		},
		HttpCode: http.StatusBadRequest,
	},
	SCHEDULEOVERLAYOPEN: {
		ResponseErr: ResponseErr{
			Message: "Agendamento fora do horário de expediente.",
			Code:    SCHEDULEOVERLAYOPEN,
		},
		HttpCode: http.StatusBadRequest,
	},
	WRONGVALUE: {
		ResponseErr: ResponseErr{
			Message: "Valor enviado inconsistente.",
			Code:    WRONGVALUE,
		},
		HttpCode: http.StatusBadRequest,
	},
	BALANCEINSUFFICIENT: {
		ResponseErr: ResponseErr{
			Message: "Saldo insuficiente.",
			Code:    BALANCEINSUFFICIENT,
		},
		HttpCode: http.StatusBadRequest,
	},
	PAYMENTPASSWORD: {
		ResponseErr: ResponseErr{
			Message: "Senha inválida.",
			Code:    PAYMENTPASSWORD,
		},
		HttpCode: http.StatusBadRequest,
	},
}
