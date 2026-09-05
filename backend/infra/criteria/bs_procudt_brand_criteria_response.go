package criteria

import "bitbucket.org/lyndus/backend/global/aggregate"

type BSProductBrand struct{
	CResponse
	Items []aggregate.BSProductBrandAggregate `json:"items"`
}

