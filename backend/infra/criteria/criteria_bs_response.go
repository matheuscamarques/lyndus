package criteria

import (
	"bitbucket.org/lyndus/backend/global/aggregate"
)

type BSResponse struct{
	CResponse
	Items []aggregate.BSAggregate `json:"items"`
}
