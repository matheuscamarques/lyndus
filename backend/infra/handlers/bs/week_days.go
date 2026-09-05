package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"net/http"
)



func GetWeekDays(w http.ResponseWriter, r *http.Request) {

	weekDays, err := services.WeekDayService.GetWeekDays()
	if err != nil {
		config.ResponsePerErr(w, err, config.NOTFOUND)
	}
	config.JSONResponse(weekDays, http.StatusOK, w)
}



