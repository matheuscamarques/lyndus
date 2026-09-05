package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)


type Scheduling struct {
	ID            int            `json:"id" db:"id"`
	BSID          int            `json:"-" db:"bs_id"`
	BSBmID        int            `json:"bmID" db:"bs_bm_id"`
	BSBmName      string         `json:"bmName"db:"bs_bm_name"`
	BSServiceID   int            `json:"serviceID" db:"bs_service_id"`
	BSServiceName string         `json:"serviceName" db:"bs_service_name"`
	BSPersonID    int            `json:"personID" db:"bs_person_id"`
	BSPersonName  string         `json:"personName" db:"bs_person_name"`
	StatusID      int            `json:"statusID" db:"schedule_status_id"`
	StartTime     types.DateTime `json:"startTime" db:"start_time"`
	EndTime       types.DateTime `json:"endTime" db:"end_time"`
	AppUser       bool           `json:"appUser" db:"app_user"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type SchedulingBAsic struct {
	ID          int            `json:"id" db:"id"`
	BSID        int            `json:"-" db:"bs_id"`
	BSBmID      int            `json:"bmID" db:"bs_bm_id"`
	BSServiceID int            `json:"serviceID" db:"bs_service_id"`
	BSPersonID  int            `json:"personID" db:"bs_person_id"`
	StatusID    int            `json:"statusID" db:"schedule_status_id"`
	StartTime   types.DateTime `json:"startTime" db:"start_time"`
	EndTime     types.DateTime `json:"endTime" db:"end_time"`
	AppUser     bool           `json:"appUser" db:"app_user"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
}

type RequestScheduling struct {
	BsID      int    `json:"bsID"`
	ServiceID int    `json:"serviceID"`
	BmID      int    `json:"bmID"`
	Date      string `json:"date"`
	Hour      string `json:"hour"`
}
