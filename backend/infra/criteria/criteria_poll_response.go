package criteria

import (
	"bitbucket.org/lyndus/backend/global/basic"
)

type CPollResponse struct {
	CResponse
	Items []basic.Poll `json:"items"`
}
