package aggregate

import (
	entity2 "bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/global/entity"
	"bitbucket.org/lyndus/backend/infra/types"
)

type BSAggregate struct {
	ID             int        `json:"id" db:"id"`
	CNPJ           types.CNPJ `json:"cnpj,omitempty" db:"cnpj"`
	entity.BS      `db:"bs"`
	entity.Company `db:"company"`
	//Latitude       *float64            `json:"latitude" db:"lat"`
	//Longitude      *float64            `json:"longitude" db:"lon"`
	Categories     []entity2.Category `json:"categories"`
	WeekDays       []WeekDay          `db:"week_days" json:"weekDays"`
}
