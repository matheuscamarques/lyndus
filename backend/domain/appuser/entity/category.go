package entity

import "bitbucket.org/lyndus/backend/global/aggregate"

type Category struct {
	ID   int    `json:"id" db:"id"`
	Desc string `json:"desc" db:"desc"`
	BS   []aggregate.BSAggregate   `json:"bs,omitempty" db:"bs"`
}

