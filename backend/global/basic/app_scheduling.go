package basic

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type  AppScheduling struct {
	ID            int            `json:"id" db:"id"`
	BsID          int            `json:"bsID" db:"bs_id"`
	BsName        string         `json:"bsName" db:"fantasy_name"`
	BSBmID        int            `json:"bmID" db:"bs_bm_id"`
	BSBmName      string         `json:"bmName"db:"bs_bm_name"`
	BSServiceID   int            `json:"serviceID" db:"bs_service_id"`
	BSServiceName string         `json:"serviceName" db:"bs_service_name"`
	Status        string         `json:"status" db:"status"`
	StatusID      int            `json:"statusID" db:"schedule_status_id"`
	StartTime     types.DateTime `json:"startTime" db:"start_time"`
	EndTime       types.DateTime `json:"endTime" db:"end_time"`
	AppUser       bool           `json:"appUser" db:"app_user"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}