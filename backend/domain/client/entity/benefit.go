package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
	"time"
)

type BenefitCreate struct {
	Desc         string `json:"desc" db:"desc"`
	AllEmployees bool   `json:"allEmployees" db:"all_employees"`
	Categories   []int  `json:"categories" db:"categories"`
}

type Benefit struct {
	ID              int             `json:"id" db:"id"`
	ClientID        int             `json:"clientID" db:"client_id,omitempty"`
	BenefitStatusID int             `json:"benefitStatusID" db:"benefit_status_id,omitempty"`
	Desc            string          `json:"desc" db:"desc"`
	Value           decimal.Decimal `json:"value" db:"value"`
	CreatedAt       time.Time       `json:"-" db:"created_at,omitempty"`
}

type BenefitUser struct {
	ID               int              `json:"ID" db:"id"`
	ClientID         int              `json:"clientID" db:"client_id"`
	BenefitID        int              `json:"benefitID" db:"benefit_id"`
	AppUserID        int              `json:"app_userID" db:"app_user_id"`
	EmployeeID       int              `json:"employeeID" db:"employee_id"`
	Value            decimal.Decimal  `json:"valID" db:"value"`
	AdditionalValue  *decimal.Decimal `json:"additionalValue" db:"additional_value"`
	AdditionalReason *string          `json:"additionalReason" db:"additional_reason"`
}

type BenefitUserUpdate struct {
	ID     int             `json:"id" db:"id"`
	Value  decimal.Decimal `json:"additionalValue" db:"value"`
	Reason string          `json:"additionalReason" db:"reason"`
}

type BenefitUserResp struct {
	ID   int       `json:"id" db:"id"`
	Name string    `json:"name"db:"name"`
	CPF  types.CPF `json:"cpf" db:"cpf"`

	Categories       []string         `json:"categories" db:"categories"`
	Value            decimal.Decimal  `json:"value" db:"value"`
	AdditionalValue  *decimal.Decimal `json:"additionalValue" db:"additional_value"`
	AdditionalReason *string          `json:"additionalReason" db:"additional_reason"`
}

type BenefitUserCategory struct {
	ID            int             `json:"id" db:"id"`
	BenefitUserID int             `json:"benefitUserID_user_id"db:"benefit_user_id"`
	CategoryID    int             `json:"categoryIDry_id"db:"category_id"`
	Name          string          `json:"name" db:"name"`
	Value         decimal.Decimal `json:"value" db:"value"`
	CreatedAt     time.Time       `json:"created_at"db:"created_at"`
}

type BenefitResp struct {
	ID              int             `json:"id" db:"id"`
	BenefitStatus   string          `json:"benefitStatus" db:"benefit_status"`
	BenefitStatusID int             `json:"benefitStatusID" db:"benefit_status_id"`
	Desc            string          `json:"desc" db:"desc"`
	Value           decimal.Decimal `json:"value" db:"value"`
	Date            time.Time       `json:"date" db:"date"`
}

type BenefitRespFull struct {
	ID              int               `json:"id" db:"id"`
	BenefitStatus   string            `json:"benefitStatus" db:"benefit_status"`
	BenefitStatusID int               `json:"benefitStatusID" db:"benefit_status_id"`
	Desc            string            `json:"desc" db:"desc"`
	Value           decimal.Decimal   `json:"value" db:"value"`
	Date            time.Time         `json:"date" db:"date"`
	BenefitUsers    []BenefitUserResp `json:"benefitUsers" db:"benefit_users"`
}
