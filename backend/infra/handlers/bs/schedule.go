package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/domain/constants"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/types"
	"bitbucket.org/lyndus/backend/infra/utils"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math"
	"strconv"

	"net/http"
	"time"
)

func GetSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.SCHEDULE, bs.VIEW)

	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	locStr, err := services.CompanyService.GetLocation(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	now := types.DateTime{}
	err = now.SetNow(locStr)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	hoje := now.Time
	dateSTR := chi.URLParam(r, "date")

	if dateSTR != "" {
		hoje, err = time.Parse(utils.LAYOUTDATE, dateSTR)
		if err != nil {
			config.ResponsePerErr(w, err, config.INVALIDREQUEST)
			return
		}
	}

	wd := int(hoje.Weekday() + 1)
	//start end bs bs_service current day
	weekDays, err := services.WeekDayService.Get(bsID, wd)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var todayScheduling entity.TodayScheduling
	n := len(weekDays)
	if n == 0 {
		config.JSONResponse(todayScheduling, http.StatusOK, w)
		return
	}

	todayScheduling.Start.Time = weekDays[0].StartTime.Add(-time.Hour)
	todayScheduling.Interval = int(math.Ceil(weekDays[n-1].EndTime.Sub(todayScheduling.Start.Time).Hours())) + 1

	//trazer os bms que só trabalham no dia corrente.
	todayScheduling.BMs, err = services.BmService.GetBMsByWeekDay(bsID, wd)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//todo rotas para iniciar atendimento, adicionar serviço, adicionar produto, finalizar, relizar pagamento, etc.

	todayScheduling.Scheduling, err = services.Schedule.GetBSScheduling(bsID, hoje)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	//tratamento para agendamentos que passaram do horário
	//mudança de status do agendamento
	for k := range todayScheduling.Scheduling {
		if todayScheduling.Scheduling[k].BSOrderID != nil {
			var order entity.Order
			order, err = services.OrderService.GetOrderByID(bsID, *todayScheduling.Scheduling[k].BSOrderID)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
			if order.ID == 0 {
				config.ResponsePerErr(w, err, config.NOTFOUND)
				return
			}

			order.Items, err = services.OrderService.GetOrdersItemsFull(order.ID)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
			todayScheduling.Scheduling[k].Order = &order

		} else {

			var service entity.Service
			service, err = services.BS.GetService(bsID, todayScheduling.Scheduling[k].BSServiceID)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
			if service.ID == 0 {
				config.ResponsePerErr(w, err, config.INVALIDREQUEST)
				return
			}

			var order entity.Order
			order.BsID = bsID
			order.BSPersonID = todayScheduling.Scheduling[k].BSPersonID
			order.Value = service.Value
			order.Discount = decimal.Zero
			order.TotalValue = service.Value
			order.Obs = ""
			order.StatusID = constants.OrderStatusOpen

			order.ID, err = services.OrderService.CreateOrder(order)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}

			var item entity.Item
			item.BSOrderID = order.ID
			item.Name = service.Name
			item.Value = service.Value
			item.BSServiceID = &service.ID
			item.BSBMID = todayScheduling.Scheduling[k].BSBmID
			item.Quantity = 1
			item.TotalValue = service.Value
			item.ID, err = services.ItemService.CreateOrderItem(item)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
			order.Items = append(order.Items, item)
			todayScheduling.Scheduling[k].Order = &order

		}
		if todayScheduling.Scheduling[k].StatusID == bs.Scheduled && now.After(todayScheduling.Scheduling[k].StartTime.Time) {
			err = services.Schedule.ChangeStatusSchedule(bsID, todayScheduling.Scheduling[k].ID, bs.Absent)
			if err != nil {
				logger.Error("Erro mudar status agenda", zap.String("err", err.Error()))
			}
			todayScheduling.Scheduling[k].StatusID = bs.Absent
		}
	}
	if todayScheduling.Scheduling == nil {
		todayScheduling.Scheduling = make([]entity.Scheduling, 0)
	}

	config.JSONResponse(todayScheduling, http.StatusOK, w)
	return
}

func NewScheduling(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.SCHEDULE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var scheduling entity.Scheduling
	err = json.NewDecoder(r.Body).Decode(&scheduling)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	wd := int(scheduling.StartTime.Weekday()) + 1
	weekDay, err := services.WeekDayService.Get(bsID, wd)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//testar se bs_bm exite
	bm, err := services.BmService.GetBM(bsID, scheduling.BSBmID, constants.BMStatusActive)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bm.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	// find bs_service by bs_bm and bs_service id
	bmService, err := services.BmService.GetBMService(bm.ID, scheduling.BSServiceID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bmService.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	// Start e End Time em relação a duração do serviço mediante ao bmxserviço escolhido.
	scheduling.EndTime.Time = scheduling.StartTime.Add(time.Duration(bmService.Duration) * time.Minute)

	start := time.Date(0, 1, 1, scheduling.StartTime.Hour(), scheduling.StartTime.Minute(), 0, 0, time.UTC)
	end := time.Date(0, 1, 1, scheduling.EndTime.Hour(), scheduling.EndTime.Minute(), 0, 0, time.UTC)

	locStr, err := services.CompanyService.GetLocation(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var now types.DateTime
	err = now.SetNow(locStr)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if scheduling.StartTime.Before(now.Time) {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	cntrl := true
	for i := range weekDay {
		if start.After(weekDay[i].StartTime.Add(time.Minute*-1)) && end.Before(weekDay[i].EndTime.Add(time.Minute+1)) {
			cntrl = false
			break
		}
	}

	if cntrl {
		config.ResponsePerErr(w, err, config.SCHEDULEOVERLAYOPEN)
		return
	}

	ids, err := services.Schedule.GetBSScheduledTime(bsID, scheduling.BSBmID, scheduling.StartTime, scheduling.EndTime)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if len(ids) > 0 {
		config.ResponsePerErr(w, err, config.SCHEDULEOVERLAY)
		return
	}
	//}else{
	//	config.ResponsePerErr(w, err, config.SCHEDULEOVERLAY)
	//	return
	//}

	// todo
	// [X] testar se bs_bm exite x serviço existe
	// [X] bs_person existe
	// [X] start e end time em relação a duração do serviço mediante ao bmxserviço escolhido.
	// [X] register date é a data do agendamento, mandado pelo front.
	// [X] front manda só o início, com a duração gerar o endtime
	// [X] checar sobreposição de horários
	// [] agendar um break para o bs_bm(agendamento sem serviço e personid)?
	// verificar se está no período de atendimento do bs e do bs_bm.

	// find bs_person by id
	person, err := services.PersonService.GetPersonID(bsID, scheduling.BSPersonID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if person.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	scheduling.BSID = bsID
	scheduling.ID, err = services.Schedule.CreateScheduling(scheduling)

	if err != nil || scheduling.ID == 0 {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	resp := response.ID{
		ID: scheduling.ID,
	}
	config.JSONResponse(resp, http.StatusOK, w)

}

func SchedulingAttendance(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.SCHEDULE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	scheduleID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	schedule, err := services.Schedule.GetBSSchedule(bsID, scheduleID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var service entity.Service
	service, err = services.BS.GetService(bsID, schedule.BSServiceID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if service.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var order entity.Order
	order.BsID = bsID
	order.BSPersonID = schedule.BSPersonID
	order.Value = service.Value
	order.Discount = decimal.Zero
	order.TotalValue = service.Value
	order.Obs = ""
	order.StatusID = constants.OrderStatusOpen

	order.ID, err = services.OrderService.CreateOrder(order)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var item entity.Item
	item.BSOrderID = order.ID
	item.Name = service.Name
	item.Value = service.Value
	item.BSServiceID = &service.ID
	item.BSBMID = schedule.BSBmID
	item.Quantity = 1
	item.TotalValue = service.Value
	item.ID, err = services.ItemService.CreateOrderItem(item)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	// todo verificar se já tem uma order do mesmo dia aberta para o mesmo cliente com agendamento em andamento.

	// todo
	// buscar dados do agendamento pelo id
	// personid, serviço, vm etc.
	// criar uma order
	// adicionar items
	// mudar status do agendamento
	//

	//bs_schedule, err := services.Schedule.GetBSSchedule(bsID, scheduleID)
	//
	// adicionar itens - produtos, serviços e insumos a uma ordem
	// remove ritens -
	// finalizar ordem
	// processar pagamento - pagamento lyndus.
	// sempre ajustando o status do serviço.
	//
	//

	err = services.Schedule.ChangeStatusSchedule(bsID, scheduleID, bs.Attendance)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = services.Schedule.SetOrderIDSchedule(bsID, scheduleID, order.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func SchedulingCharge(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.SCHEDULE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	scheduleID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var paymentsCharge []entity.PaymentCharge
	err = json.NewDecoder(r.Body).Decode(&paymentsCharge)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	schedule, err := services.Schedule.GetBSSchedule(bsID, scheduleID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if schedule.BSOrderID == nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if schedule.StatusID != bs.Attendance {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var order entity.Order
	order, err = services.OrderService.GetOrderByID(bsID, *schedule.BSOrderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if order.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if order.StatusID == constants.OrderStatusCanceled {
		config.ResponsePerErr(w, err, config.ORDERCANCELED)
		return
	} else if order.StatusID == constants.OrderStatusExpired {
		config.ResponsePerErr(w, err, config.ORDEREXPIRED)
		return
	}

	var value decimal.Decimal
	for _, payment := range paymentsCharge {
		value = value.Add(payment.Value)
	}

	////log.Println(value, order.TotalValue)
	if value != order.TotalValue {
		config.ResponsePerErr(w, err, config.ORDERINVALIDVALUE)
		return
	}
	lyndusPay := false
	for _, payment := range paymentsCharge {
		payment.BSOrderID = *schedule.BSOrderID
		payment.ID, err = services.PaymentService.CreateOrderPayment(payment)
		if payment.PaymentTypeID == bs.LyndusPayment {
			lyndusPay = true
			// todo pagamento pelo app
			paymentUUID, err := uuid.NewUUID()
			if err != nil {
				logger.Error("Err on generate UUID PAYMENT:", zap.String("ERROR:", err.Error()))
			}
			balanceReceivable := entity.BalanceReceivable{
				BsID:             bsID,
				BSPersonID:       order.BSPersonID,
				BSOrderPaymentID: payment.ID,
				PaymentStatusID:  bs.PaymentStatusWaitingPayment,
				PaymentUUID:      paymentUUID,
				Value:            payment.Value,
			}
			_, err = services.PaymentService.CreateOrderReceivable(balanceReceivable)
			if err != nil {
				logger.Error("Err on generate UUID PAYMENT:", zap.String("ERROR:", err.Error()))
			}
		}
	}
	if lyndusPay {
		schedule.StatusID = bs.WaitPayment
		order.StatusID = constants.OrderStatusWaitingPayment
	} else {
		schedule.StatusID = bs.Finished
		order.StatusID = constants.OrderStatusPaid
	}

	err = services.OrderService.UpdateOrderStatus(bsID, order.ID, order.StatusID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = services.Schedule.ChangeStatusSchedule(bsID, scheduleID, schedule.StatusID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

// CancelScheduling
// Todo cancelar agendamento revisar. cancelamento não tem restrição.
//CancelScheduling cancel scheduling
func CancelScheduling(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.SCHEDULE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	scheduleID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	schedule, err := services.Schedule.GetBSSchedule(bsID, scheduleID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if schedule.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	//todo validar regras, tratamento e retorno
	if schedule.StatusID == bs.WaitPayment || schedule.StatusID == bs.Finished || schedule.StatusID == bs.Canceled {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.Schedule.ChangeStatusSchedule(bsID, scheduleID, constants.SCHEDULECANCELED)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if schedule.BSOrderID != nil {

		err = services.OrderService.UpdateOrderStatus(bsID, *schedule.BSOrderID, constants.OrderStatusCanceled)

		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		receivables, err := services.PaymentService.ListOrderReceivable(bsID, *schedule.BSOrderID)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		for k := range receivables {
			err = services.PaymentService.UpdateOrderReceivableStatus(bsID, receivables[k].ID, bs.PaymentStatusCanceled)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}
