package aggregate

import "bitbucket.org/lyndus/backend/infra/types"

type WeekDay struct {
	ID        int            `json:"id" db:"id"`
	Name      string         `json:"name,omitempty" db:"name"`
	StartTime types.TimeHHMM `json:"startTime" db:"start_time"`
	EndTime   types.TimeHHMM `json:"endTime" db:"end_time"`
}