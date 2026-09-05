package bonder

import (
	"github.com/go-chi/chi"
)

func BSRouter(r chi.Router) {
	//r.Use(auth.ValidateTokenMiddleware)
	r.Get("/categories", BsGetCategoriesHandler)
	r.Post("/", BsRegisterHandler)
	r.Put("/", BsUpdateHandler)
	r.Get("/", BsGetAllHandler)
	r.Get("/{id}", BsGetByIdHandler)
	r.Put("/{id}/active", BsActivateHandler)
	r.Put("/{id}/inactive", BsInativateHandler)
	r.Delete("/{id}", BsDeleteHandler)
	r.Post("/{id}/profile/logo", BsSaveProfileLogo)
	r.Get("/{id}/profile/logo", BsGetProfileLogo)
	//{{host_homolog_bs}}profile/logo
	r.Get("/payments", BsGetAllPayments)
	r.Get("/payments/{id}", BsGetByIdPayment)
	r.Post("/payments/{id}", BSConfirmPayment)
	r.Get("/payments/{id}/file", BSPaymentConfirmFile)

}
func BMRouter(r chi.Router) {
	r.Get("/financial/withdraw", ListWithdrawalRequest)
	r.Get("/financial/withdraw/{withdrawalID}", WithdrawalRequestByID)
	r.Post("/financial/withdraw/{withdrawalID}", BMConfirmPayment)
	r.Get("/financial/withdraw/{withdrawalID}/file", BMPaymentConfirmFile)
}

func ClientRouter(r chi.Router) {
	//r.Use(auth.ValidateTokenMiddleware)
	r.Post("/", ClientRegisterHandler)
	r.Get("/", ClientGetAllHandler)
	r.Get("/{id}", ClientGetByIdHandler)
	r.Put("/{id}/active", ClientActiveHandler)
	r.Put("/{id}/inactive", ClientInactiveHandler)
	r.Delete("/{id}", ClientDeleteHandler)
	r.Put("/", ClientUpdateHandler)

	r.Get("/benefit", BenefitsGetAllHandler)
	r.Get("/benefit/{id}", BenefitGetHandler)
	r.Post("/benefit/{id}/sendTicket", BenefitSendTicketHandler)
	r.Put("/benefit/{id}/confirmPayment", BenefitConfirmPaymentHandler)

}

func AppUserRouter(r chi.Router) {
	r.Get("/", AppuserGetAll)
	r.Get("/{id}", AppUserGetByID)
	r.Post("/bonus", AppUserUpdateLyndusBoxBonus)
}

func MasterRouter(r chi.Router) {
	r.Post("/client/change-password", MasterUserClientChangePassword)
	r.Post("/bs/change-password", MasterUserBSChangePassword)
}
