package appuser

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/services"
	"bitbucket.org/lyndus/backend/domain/constants"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/types"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// todo
//
// buscar bs_bm
// buscar serviços do bs_bm
// listar horários livres
// agendamento
// cancelamento
// listar agendamento do app bs_user, mostrar agendamento.
//

func GetSchedules(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	dateStr := r.URL.Query().Get("date")

	var date types.Date
	if dateStr != "" {
		err := date.ParseISO(dateStr)
		if err != nil {
			config.ResponsePerErr(w, err, config.INVALIDREQUEST)
			return
		}
	}

	if page == 0 {
		page = 1
	}
	if itemsPerPage == 0 || itemsPerPage > 100 {
		itemsPerPage = 10
	}

	var err error
	var response criteria.AppScheduleResponse
	response.Items, response.TotalPages, err = services.Schedule.GetSchedulingsByAppUserID(id, page, itemsPerPage, date)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	for k := range response.Items {
		response.Items[k].Status = constants.SchedulesStatus[response.Items[k].StatusID]
		response.Items[k].BsName, err = services.BSService.GetBSName(response.Items[k].BsID)
		if err != nil {
			logger.Error("Erro ao listar agendamento bs_bm", zap.String("err: ", err.Error()))
		}
	}

	response.ActivePage = page

	config.JSONResponse(response, http.StatusOK, w)
	return
}

func GetBsServicesBmHours(w http.ResponseWriter, r *http.Request) {

	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	serviceID, err := strconv.Atoi(chi.URLParam(r, "serviceID"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	bmID, err := strconv.Atoi(chi.URLParam(r, "bmID"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if bsID == 0 || serviceID == 0 || bmID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	month, _ := strconv.Atoi(r.URL.Query().Get("month"))

	//todo buscar bsid x bmid ativo
	bm, err := services.BM.GetBMByIDAndDuration(bsID, bmID, serviceID)
	if err == sql.ErrNoRows || bm.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//
	// buscar horário funcionamento bs ignorado por hora,
	// presumisse que o controle no cadastro do bs_bm seja suficiente
	// bsWeedDays, err := services.BSService.GetBSWeekDays(bsID)
	// if err != nil{
	//	config.ResponsePerErr(w, err, config.INVALIDREQUEST)
	//	return
	// }

	// buscar horario trabalho bs_bm
	bmWeedDays, err := services.BM.GetBMWeekDays(bmID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	locStr, err := services.BSService.GetLocation(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	// gerar as horas de trabalho do bs_bm pelo step time dele
	step := time.Minute * time.Duration(bm.ServiceStepTime)
	// gerar as horas do dia today
	//gerar as horas dos próximos dias.
	var today, lastDay, tempControl1, tempControl2, tempControl3 types.Date

	location, err := time.LoadLocation(locStr)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	now := time.Now().In(location)
	if month > 0 && month <= 12 {
		//tratamento para ignorar fuso horário
		now = time.Date(now.Year(), time.Month(month), 1, now.Hour(), now.Minute(), 0, 0, &time.Location{})
	} else {
		//tratamento para ignorar fuso horário
		now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, &time.Location{})
	}
	//tratamento para zerar hh:mm:ss e fuso horário
	today.Time = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, &time.Location{})

	var currentWD int
	var h types.TimeHHMM
	dias := make(map[string][]types.TimeHHMM)

	// comparar dias do bs x dias do bs_bm
	// todos os agendamentos do bs + bs_bm  no período  "dias"
	// remover as horas que estão contidas no intervalo de agendamento.
	// verificar os agendamentos para excluir as horas geradas/disponíveis
	// o tempo do serviço escolhido tem que encaixar na janela de tempo disponível

	lastDay.Time = today.AddDate(0, 1, -today.Day())

	dayLimit := int(lastDay.Sub(today.Time).Hours()/24) + 1
	//log.Println(dayLimit, lastDay, today, lastDay.Sub(today.Time).Hours()/24)

	//dayLimit := 5
	//lastDay.Time = today.AddDate(0, 0, dayLimit)

	scheduling, err := services.Schedule.GetBSBMScheduling(bsID, bm.ID, today.Time, lastDay.Time)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//day offs bs_bm
	bmDaysOff, err := services.BM.GetBMDaysOff(bsID, bmID, today.Time, lastDay.Time)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//day offs bs
	bsDaysOff, err := services.BSService.GetBSDaysOff(bsID, today.Time, lastDay.Time)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bmDaysOff = append(bmDaysOff, bsDaysOff...)

	var hours []types.TimeHHMM
	serviceDuration := time.Duration(bm.ServiceDuration) * time.Minute
	tempControl1 = today

	for i := 0; i < dayLimit; i++ {
		tempControl1 = today
		currentWD = int(today.Weekday()) + 1
		hours = []types.TimeHHMM{}
		for j := range bmWeedDays {
			if bmWeedDays[j].WeekDayID == currentWD {
				h = bmWeedDays[j].StartTime
			WeekDaysFor:
				for h.Before(bmWeedDays[j].EndTime.Time) {
					//a partir da hora atual do dia atual
					tempControl2.Time = tempControl1.Add(time.Duration(h.Hour())*time.Hour + time.Duration(h.Minute())*time.Minute)
					if tempControl2.Before(now) {
						h.Time = h.Add(step)
						continue
					}
					tempControl3.Time = tempControl2.Add(serviceDuration)

					// tempo disponível mais serviço não pode ultrapassar final do expediente.
					if h.Add(serviceDuration).After(bmWeedDays[j].EndTime.Time) {
						h.Time = h.Add(step)
						continue
					}
					// remove horas agendadas
					for k := range scheduling {
						if ((tempControl2.After(scheduling[k].StartTime.Time) || tempControl2.Equal(scheduling[k].StartTime.Time)) &&
							tempControl2.Before(scheduling[k].EndTime.Time)) ||
							(tempControl3.After(scheduling[k].StartTime.Time) && tempControl3.Before(scheduling[k].EndTime.Time)) {
							h.Time = h.Add(step)
							continue WeekDaysFor
						}
					}
					// remove por day offs
					for k := range bmDaysOff {
						if ((tempControl2.After(bmDaysOff[k].StartDate.Time) ||
							tempControl2.Equal(bmDaysOff[k].StartDate.Time)) &&
							tempControl2.Before(bmDaysOff[k].EndDate.Time)) ||
							((tempControl3.After(bmDaysOff[k].StartDate.Time)) &&
								tempControl3.Before(bmDaysOff[k].EndDate.Time)) {
							h.Time = h.Add(step)
							continue WeekDaysFor
						}
					}
					hours = append(hours, h)
					h.Time = h.Add(step)
				}
			}
		}
		if len(hours) == 0 {
			dayLimit += 1
		} else {
			dias[today.String()] = hours
		}
		today.Time = today.AddDate(0, 0, 1)

	}

	config.JSONResponse(dias, http.StatusOK, w)
}

// GetBsServicesByDate listar horários by bs/serviço/data
func GetBsServicesByDate(w http.ResponseWriter, r *http.Request) {
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	serviceID, err := strconv.Atoi(chi.URLParam(r, "serviceID"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if bsID == 0 || serviceID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	pDate := chi.URLParam(r, "date")
	var date, today types.Date
	err = date.ParseISO(pDate)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	today.SetNow()
	if date.Before(today.Time) {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	weekDay := int(date.Weekday()) + 1
	// listar bms que trabalham no dia atual,
	bms, err := services.BM.GetBMsFullByWeekDay(bsID, serviceID, weekDay)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	//day offs bs
	bsDaysOff, err := services.BSService.GetBSDaysOff(bsID, today.Time, today.Time)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	locStr, err := services.BSService.GetLocation(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	location, err := time.LoadLocation(locStr)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	now := time.Now().In(location)
	//tratamento para ignorar fuso horário
	now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, &time.Location{})
	var tempControl1, tempControl2, tempControl3 types.Date
	hours := make(map[int64]types.TimeHHMM)
	var h types.TimeHHMM

	tempControl1 = today
	for i := range bms {
		h = *bms[i].StartTime
		scheduling, err := services.Schedule.GetBSBMScheduling(bsID, bms[i].ID, today.Time, today.Time)
		if err != nil {
			logger.Error("Erro ao listar agendamento bs_bm", zap.String("err: ", err.Error()))
		}
		bmDaysOff, err := services.BM.GetBMDaysOff(bsID, bms[i].ID, today.Time, today.Time)
		if err != nil {
			logger.Error("Erro ao listar day offs BM", zap.String("err: ", err.Error()))
		}
		bmDaysOff = append(bmDaysOff, bsDaysOff...)

	WeekDaysFor:
		for h.Before(bms[i].EndTime.Time) {
			//a partir da hora atual do dia atual
			tempControl2.Time = tempControl1.Add(time.Duration(h.Hour())*time.Hour + time.Duration(h.Minute())*time.Minute)
			if tempControl2.Before(now) && today.Equal(date.Time) {
				h.Time = h.Add(bms[i].StepTimeMinutes())
				continue
			}

			tempControl3.Time = tempControl2.Add(bms[i].ServiceDurationMinutes())

			// tempo disponível mais serviço não pode ultrapassar final do expediente.
			if h.Add(bms[i].ServiceDurationMinutes()).After(bms[i].EndTime.Time) {
				h.Time = h.Add(bms[i].StepTimeMinutes())
				continue
			}

			// remove horas agendadas
			for k := range scheduling {
				if ((tempControl2.After(scheduling[k].StartTime.Time) || tempControl2.Equal(scheduling[k].StartTime.Time)) &&
					tempControl2.Before(scheduling[k].EndTime.Time)) ||
					(tempControl3.After(scheduling[k].StartTime.Time) && tempControl3.Before(scheduling[k].EndTime.Time)) {
					h.Time = h.Add(bms[i].StepTimeMinutes())
					continue WeekDaysFor
				}
			}

			// remove por day offs
			for k := range bmDaysOff {
				if ((tempControl2.After(bmDaysOff[k].StartDate.Time) ||
					tempControl2.Equal(bmDaysOff[k].StartDate.Time)) &&
					tempControl2.Before(bmDaysOff[k].EndDate.Time)) ||
					((tempControl3.After(bmDaysOff[k].StartDate.Time)) &&
						tempControl3.Before(bmDaysOff[k].EndDate.Time)) {
					h.Time = h.Add(bms[i].StepTimeMinutes())
					continue WeekDaysFor
				}
			}
			hours[h.Unix()] = h
			h.Time = h.Add(bms[i].StepTimeMinutes())
		}
	}

	var respHours []types.TimeHHMM
	for _, key := range hours {
		respHours = append(respHours, key)
	}
	sort.Slice(respHours, func(i, j int) bool { return respHours[i].UnixNano() < respHours[j].UnixNano() })

	config.JSONResponse(respHours, http.StatusOK, w)
}

func GetBsBmServiceByDateAndHour(w http.ResponseWriter, r *http.Request) {
	//id serviceID date hour
	bsID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	serviceID, err := strconv.Atoi(chi.URLParam(r, "serviceID"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	sDate := chi.URLParam(r, "date")
	sHour := chi.URLParam(r, "hour")

	var date types.Date
	err = date.ParseISO(sDate)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	var hour types.TimeHHMM
	hour.ParseISO(sHour)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	fullDate := time.Date(date.Year(), date.Month(), date.Day(), hour.Hour(), hour.Minute(), 0, 0, &time.Location{})

	locStr, err := services.BSService.GetLocation(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	location, err := time.LoadLocation(locStr)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	now := time.Now().In(location)

	//tratamento para ignorar fuso horário
	now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, &time.Location{})

	if fullDate.Before(now) {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	wd := int(date.Weekday()) + 1

	bms, err := services.BM.GetBMsFullByWeekDay(bsID, serviceID, wd)
	if len(bms) <= 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//day offs bs_bm
	bsDaysOff, err := services.BSService.GetBSDaysOff(bsID, date.Time, date.Time)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	offDay := false
	for i := range bsDaysOff {
		if (fullDate.After(bsDaysOff[i].StartDate.Time) ||
			fullDate.Equal(bsDaysOff[i].StartDate.Time)) &&
			fullDate.Before(bsDaysOff[i].EndDate.Time) {
			offDay = true
			break
		}
	}
	if offDay {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	var step time.Duration
	var filterBms []entity.BM
	var tempControl1, tempControl2, tempControl3 types.Date

	tempControl1 = date
	for _, bm := range bms {

		// gerar as horas de trabalho do bs_bm pelo step time dele
		//comparar dias do bs x dias do bs_bm
		step = time.Minute * time.Duration(bm.ServiceStepTime)
		// gerar as horas do dia today
		//gerar as horas dos próximos dias.

		var h types.TimeHHMM
		// todos os agendamentos do bs + bs_bm  no período  "dias"
		// remover as horas que estão contidas no intervalo de agendamento.
		// verificar os agendamentos para excluir as horas geradas/disponíveis
		// o tempo do serviço escolhido tem que encaixar na janela de tempo disponível

		scheduling, err := services.Schedule.GetBSBMScheduling(bsID, bm.ID, date.Time, date.Time)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		bmDaysOff, err := services.BM.GetBMDaysOff(bsID, bm.ID, date.Time, date.Time)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		serviceDuration := bm.ServiceDurationMinutes()

		h = *bm.StartTime
	WeekDaysFor:
		for h.Before(bm.EndTime.Time) {
			//a partir da hora atual do dia atual
			tempControl2.Time = tempControl1.Add(time.Duration(h.Hour())*time.Hour + time.Duration(h.Minute())*time.Minute)
			tempControl3.Time = tempControl2.Add(serviceDuration)

			// tempo disponível mais serviço não pode ultrapassar final do expediente.
			if h.Add(serviceDuration).After(bm.EndTime.Time) {
				h.Time = h.Add(step)
				continue
			}
			// remove horas agendadas
			for k := range scheduling {
				if ((tempControl2.After(scheduling[k].StartTime.Time) || tempControl2.Equal(scheduling[k].StartTime.Time)) &&
					tempControl2.Before(scheduling[k].EndTime.Time)) ||
					(tempControl3.After(scheduling[k].StartTime.Time) && tempControl3.Before(scheduling[k].EndTime.Time)) {
					h.Time = h.Add(step)
					continue WeekDaysFor
				}
			}
			// remove por day offs
			for k := range bmDaysOff {
				if ((tempControl2.After(bmDaysOff[k].StartDate.Time) ||
					tempControl2.Equal(bmDaysOff[k].StartDate.Time)) &&
					tempControl2.Before(bmDaysOff[k].EndDate.Time)) ||
					((tempControl3.After(bmDaysOff[k].StartDate.Time)) &&
						tempControl3.Before(bmDaysOff[k].EndDate.Time)) {
					h.Time = h.Add(step)
					continue WeekDaysFor
				}
			}
			if hour.Equal(h.Time) {
				filterBms = append(filterBms, entity.BM{
					ID:              bm.ID,
					Name:            bm.Name,
					Obs:             bm.Obs,
					ServiceDuration: bm.ServiceDuration,
				})
				break
			}
			h.Time = h.Add(step)
		}
	}
	if len(filterBms) == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	config.JSONResponse(filterBms, http.StatusOK, w)
}

//CreateSchedule salvar novo agendamento.
func CreateSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	var reqSchedule entity.RequestScheduling
	err := json.NewDecoder(r.Body).Decode(&reqSchedule)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var startTime, endTime, now types.DateTime
	var date types.Date
	var hour types.TimeHHMM

	// use date and hour to create
	err = startTime.ParseISO(reqSchedule.Date + " " + reqSchedule.Hour)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	err = date.ParseISO(reqSchedule.Date)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	err = hour.ParseISO(reqSchedule.Hour)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	wd := int(date.Weekday()) + 1
	bm, err := services.BM.GetBMFullByWeekDay(reqSchedule.BsID, reqSchedule.ServiceID, reqSchedule.BmID, wd)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bm.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	endTime.Time = startTime.Add(bm.ServiceDurationMinutes())

	bsDaysOff, err := services.BSService.GetBSDaysOff(reqSchedule.BsID, startTime.Time, endTime.Time)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var tempControl1, tempControl2, tempControl3 types.Date

	// gerar as horas de trabalho do bs_bm pelo step time dele
	//comparar dias do bs x dias do bs_bm
	step := time.Minute * time.Duration(bm.ServiceStepTime)
	// gerar as horas do dia today
	//gerar as horas dos próximos dias.

	var h types.TimeHHMM
	// todos os agendamentos do bs + bs_bm  no período  "dias"
	// remover as horas que estão contidas no intervalo de agendamento.
	// verificar os agendamentos para excluir as horas geradas/disponíveis
	// o tempo do serviço escolhido tem que encaixar na janela de tempo disponível

	schedules, err := services.Schedule.GetBSBMScheduling(reqSchedule.BsID, bm.ID, date.Time, date.Time)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	bmDaysOff, err := services.BM.GetBMDaysOff(reqSchedule.BsID, bm.ID, date.Time, date.Time)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	bmDaysOff = append(bmDaysOff, bsDaysOff...)

	serviceDuration := bm.ServiceDurationMinutes()

	locStr, err := services.BSService.GetLocation(reqSchedule.BsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = now.SetNow(locStr)
	if startTime.Before(now.Time) {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	tempControl1 = date
	freeHour := false
	h = *bm.StartTime
WeekDaysFor:
	for h.Before(bm.EndTime.Time) {
		//a partir da hora atual do dia atual
		tempControl2.Time = tempControl1.Add(time.Duration(h.Hour())*time.Hour + time.Duration(h.Minute())*time.Minute)
		tempControl3.Time = tempControl2.Add(serviceDuration)

		// tempo disponível mais serviço não pode ultrapassar final do expediente.
		if h.Add(serviceDuration).After(bm.EndTime.Time) {
			break
		}
		if hour.Before(h.Time) {
			break
		}
		// remove horas agendadas
		for k := range schedules {
			if ((tempControl2.After(schedules[k].StartTime.Time) || tempControl2.Equal(schedules[k].StartTime.Time)) &&
				tempControl2.Before(schedules[k].EndTime.Time)) ||
				(tempControl3.After(schedules[k].StartTime.Time) && tempControl3.Before(schedules[k].EndTime.Time)) {
				h.Time = h.Add(step)
				continue WeekDaysFor
			}
		}
		// remove por day offs
		for k := range bmDaysOff {
			if ((tempControl2.After(bmDaysOff[k].StartDate.Time) ||
				tempControl2.Equal(bmDaysOff[k].StartDate.Time)) &&
				tempControl2.Before(bmDaysOff[k].EndDate.Time)) ||
				((tempControl3.After(bmDaysOff[k].StartDate.Time)) &&
					tempControl3.Before(bmDaysOff[k].EndDate.Time)) {
				h.Time = h.Add(step)
				continue WeekDaysFor
			}
		}
		if hour.Equal(h.Time) {
			freeHour = true
			break
		}
		h.Time = h.Add(step)
	}

	if !freeHour {
		config.ResponsePerErr(w, err, config.SCHEDULEOVERLAY)
		return
	}

	scheduling := entity.Scheduling{
		BSBmID:      reqSchedule.BmID,
		BSID:        reqSchedule.BsID,
		BSServiceID: reqSchedule.ServiceID,
		StartTime:   startTime,
		EndTime:     endTime,
		AppUser:     true,
	}

	// buscar bs bs_person por appUser id
	bsPerson, err := services.PersonService.GetBSPersonByAppUserID(reqSchedule.BsID, id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bsPerson.ID == 0 {
		// buscar bs_person por cpf
		appUser, err := services.UserUserService.GetAppUserByID(id)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		bsPerson, err = services.PersonService.GetBSPersonByCPF(reqSchedule.BsID, appUser.CPF)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		if bsPerson.ID == 0 {
			bsPerson.CPF = appUser.CPF
			bsPerson.Birthdate = appUser.Birthdate
			bsPerson.Name = appUser.Name
			bsPerson.Phone = appUser.Phone
			bsPerson.BsID = reqSchedule.BsID
			bsPerson.AppUserID = &appUser.ID

			//criar bsPerson
			bsPerson.ID, err = services.PersonService.CreatePerson(bsPerson)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}
		if bsPerson.AppUserID != nil && *bsPerson.AppUserID == 0 {
			//update appuserid
			err = services.PersonService.UpdatePersonAppUserID(reqSchedule.BsID, bsPerson.ID, id)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}

	}
	if bsPerson.ID == 0 {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return

	}

	scheduling.BSPersonID = bsPerson.ID
	idSchedule, err := services.Schedule.CreateScheduling(scheduling)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(response.ID{ID: idSchedule}, http.StatusOK, w)

}

//DeleteSchedule salvar novo agendamento.
func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	scheduleID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	log.Println(id, scheduleID)
	schedule, err := services.Schedule.GetSchedulingAppUser(id, scheduleID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if schedule.ID == 0 || schedule.ID != scheduleID {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	err = services.Schedule.ChangeStatusSchedule(schedule.BSID, schedule.ID, constants.SCHEDULECANCELED)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
