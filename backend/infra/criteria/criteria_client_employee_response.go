package criteria

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
)

type ClientEmployeeResponse struct {
	CResponse
	Items  []entity.Employee `json:"items"`
}