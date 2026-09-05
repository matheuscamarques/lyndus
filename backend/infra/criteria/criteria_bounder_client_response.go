package criteria

import (
	"bitbucket.org/lyndus/backend/global/aggregate"
)

type BounderClientResponse struct {
	CResponse
	Items  []aggregate.ClientAggregate `json:"items"`
}
