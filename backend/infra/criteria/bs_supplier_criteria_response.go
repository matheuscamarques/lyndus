package criteria

import "bitbucket.org/lyndus/backend/global/aggregate"

type BSSupplier struct{
	CResponse
	Items []aggregate.BSSupplierAggregate `json:"items"`
}

