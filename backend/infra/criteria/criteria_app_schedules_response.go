package criteria

import (
	"bitbucket.org/lyndus/backend/global/basic"
)

type AppScheduleResponse struct {
	CResponse
	Items []basic.AppScheduling `json:"items"`
}
