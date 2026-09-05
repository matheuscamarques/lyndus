package entity

import (
	base "bitbucket.org/lyndus/backend/global/entity"
)

type Company struct {
	base.Company
	WeekDays []BMWeekDay `json:"weekDays" db:"weekDays"`
}
