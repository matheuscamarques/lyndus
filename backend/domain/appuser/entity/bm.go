package entity

import (
	"bitbucket.org/lyndus/backend/global/aggregate"
	"time"

	"bitbucket.org/lyndus/backend/infra/types"
)

type BM struct {
	ID              int             `json:"id" db:"id"`
	Name            string          `json:"name" db:"name"`
	Obs             string          `json:"obs" db:"obs"`
	ServiceDuration int             `json:"serviceDuration" db:"duration"`
	ServiceStepTime int             `json:"serviceStepTime,omitempty" db:"service_step_time"`
	StartTime       *types.TimeHHMM `json:"startTime,omitempty" db:"start_time"`
	EndTime         *types.TimeHHMM `json:"endTime,omitempty" db:"end_time"`
	Desc            string          `json:"desc,omitempty" db:"desc"`

	WeekDays []aggregate.WeekDay `json:"weekDays,omitempty" db:"week_days"`
}

func (b BM) StepTimeMinutes() time.Duration {
	return time.Minute * time.Duration(b.ServiceStepTime)
}

func (b BM) ServiceDurationMinutes() time.Duration {
	return time.Duration(b.ServiceDuration) * time.Minute
}
