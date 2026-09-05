package calendar

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

//NextDays Retorna a partir de uma data os próximos N dias
func NextDays(nDays int,days map[string][]types.TimeHHMM,now types.Date,bmWeedDays []entity.WeekDayBS,step time.Duration)  bool {
	if nDays <= 0 {
		return true
	}

	days[now.String()] = MakeHours(now,bmWeedDays,step)
	now.Time = now.AddDate(0,0,1)
	return NextDays(nDays-1,days,now,bmWeedDays,step)
}
//,bmWeedDays []entity.WeekDayBS,step time.Duration

func MakeHours(now types.Date,bmWeedDays []entity.WeekDayBS,step time.Duration) []types.TimeHHMM{
	var hours []types.TimeHHMM
	var h types.TimeHHMM

	currentWD := int(now.Weekday())+1
	for k := range bmWeedDays {
		if bmWeedDays[k].WeekDayID == currentWD {
			h = bmWeedDays[k].StartTime
			for h.Before(bmWeedDays[k].EndTime.Time) || h.Equal(bmWeedDays[k].EndTime.Time) {

				hours = append(hours, h)
				h.Time = h.Add(step)
			}
		}
	}
	return hours
}


