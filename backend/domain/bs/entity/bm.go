package entity

import (
	"github.com/shopspring/decimal"
	"log"
	"time"

	"bitbucket.org/lyndus/backend/infra/types"
)

//type BMCreate struct {
//	ID       int       `json:"id" db:"id"`
//	Name     string    `json:"name" db:"name"`
//	Document types.DOC `json:"document" db:"document"`
//	OBS      string    `json:"obs" db:"obs"`
//	Image    string    `json:"image" db:"image"`
//	// services vem como []int
//	Services []int `json:"services" db:"services"`
//
//	WeekDays []BMWeekDay `json:"weekDays" db:"week_days"`
//
//	CommissionService decimal.Decimal `json:"commisionService" db:"commission_service"`
//	CommissionProduct decimal.Decimal `json:"commisionProduct" db:"commission_product"`
//	CommissionSupply  decimal.Decimal `json:"commisionSupply"  db:"commission_supply"`
//	ServiceStepTime   decimal.Decimal `json:"serviceStepTime"  db:"service_step_time"`
//}

type BM struct {
	ID        int         `json:"id" db:"id"`
	BsID      int         `json:"bsID,omitempty" db:"bs_id"`
	Name      string      `json:"name" db:"name"`
	CPF       types.CPF   `json:"cpf" db:"cpf"`
	CNPJ      *types.CNPJ `json:"cnpj" db:"cnpj"`
	OBS       string      `json:"obs" db:"obs"`
	CreatedAt *time.Time  `json:"createdAT,omitempty" db:"created_at"`
	StatusID  int         `json:"bmStatusID,omitempty" db:"bm_status_id"`
	Image     string      `json:"image" db:"image"`
	Desc      string      `json:"desc" db:"desc"`
	// services vem como []Service
	Services          []Service       `json:"services" db:"services"`
	WeekDays          []BMWeekDay     `json:"weekDays" db:"week_days"`
	CommissionService decimal.Decimal `json:"commissionService" db:"commission_service"`
	CommissionProduct decimal.Decimal `json:"commissionProduct" db:"commission_product"`
	CommissionSupply  decimal.Decimal `json:"commissionSupply"  db:"commission_supply"`
	ServiceStepTime   int             `json:"serviceStepTime"  db:"service_step_time"`
}

type BMBasic struct {
	ID   int    `json:"id" db:"id"`
	BsID int    `json:"bsID,omitempty" db:"bs_id"`
	Name string `json:"name" db:"name"`
}

func (bm BM) ValidateServiceStepTime() bool {
	for i := range bm.Services {
		log.Println(bm.Services[0].Duration, bm.ServiceStepTime)
		if (bm.Services[i].Duration % bm.ServiceStepTime) != 0 {
			return false
		}
	}
	return true
}
