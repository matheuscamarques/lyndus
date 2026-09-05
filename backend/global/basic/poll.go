package basic

import (
	"bitbucket.org/lyndus/backend/infra/types"
)

type Poll struct {
	ID          int            `json:"id" db:"id"`
	ClientID    int 		   `json:"-" db:"client_id"`
	StartTime   types.DateTime `json:"startTime" db:"start_time"`
	EndTime     types.DateTime `json:"endTime" db:"end_time"`
	AllEmployee bool           `json:"allEmployee" db:"all_employee"`
	Title       string         `json:"title" db:"title"`
	Status      uint8            `json:"status" db:"status"`
}
