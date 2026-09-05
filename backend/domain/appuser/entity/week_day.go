package entity

import "bitbucket.org/lyndus/backend/infra/types"

type WeekDay struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type BMWeekDay struct {
	ID        int            `json:"id" db:"id"`
	Name      string         `json:"name,omitempty" db:"name"`
	StartTime types.TimeHHMM `json:"startTime" db:"start_time"`
	EndTime   types.TimeHHMM `json:"endTime" db:"end_time"`
}

type WeekDayBS struct {
	ID        int            `json:"id" db:"id"`
	BSID      int            `json:"bsID" db:"bs_id"`
	WeekDayID int            `json:"weekDayID" db:"week_day_id"`
	StartTime types.TimeHHMM `json:"startTime" db:"start_time"`
	EndTime   types.TimeHHMM `json:"endTime" db:"end_time"`
}
