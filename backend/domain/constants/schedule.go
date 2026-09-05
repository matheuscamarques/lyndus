package constants

const (
	BMStatusActive   = 1
	BMStatusInactive = 2

	//SCHEDULEACTIVATE = 1
	//SCHEDULECANCEL   = 2

	SCHEDULESCHEDULED   = 1
	SCHEDULECANCELED    = 2
	SCHEDULEATTEDANCE   = 3
	SCHEDULEWAITPAYMENT = 4
	SCHEDULEFINISHED    = 5
	SCHEDULEABSENT      = 6
)

var SchedulesStatus = map[int]string{
	SCHEDULESCHEDULED:   "Agendado",
	SCHEDULECANCELED:    "Cancelado",
	SCHEDULEATTEDANCE:   "Em atendimento",
	SCHEDULEWAITPAYMENT: "Aguardando Pagamento",
	SCHEDULEFINISHED:    "Finalizado",
	SCHEDULEABSENT:      "Ausente",
}
