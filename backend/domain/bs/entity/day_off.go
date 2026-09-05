package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type DayOff struct {
	ID        int            `json:"id" db:"id"`
	StartDate types.DateTime `json:"startDate" db:"start_date"`
	EndDate   types.DateTime `json:"endDate" db:"end_date"`
	CreatedAt time.Time      `json:"-" db:"created_at"`
}
