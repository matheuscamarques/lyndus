package response

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"github.com/shopspring/decimal"
)

type ID struct {
	ID int `json:"id"`
}

type BenefitStatus struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type Value struct {
	Value decimal.Decimal `json:"value"`
}

type BSBMResponse struct {
	ActivePage int `json:"activePage"`
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`

	Items []entity.BM `json:"items"`
}

type ListResponse struct {
	Active     bool `json:"active"`
	ActivePage int  `json:"activePage"`
	TotalPages int  `json:"totalPages"`
	TotalItems int  `json:"totalItems"`

	Items interface{} `json:"items"`
}
