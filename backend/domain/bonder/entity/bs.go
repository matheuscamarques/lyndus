package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/shopspring/decimal"
	"time"
)

type WeekDay struct {
	ID        int            `json:"id" db:"id"`
	Name      string         `json:"name,omitempty" db:"name"`
	StartTime types.TimeHHMM `json:"startTime" db:"start_time"`
	EndTime   types.TimeHHMM `json:"endTime" db:"end_time"`
}

type BS struct {
	ID                int        `json:"id" db:"id"`
	CNPJ              types.CNPJ `json:"cnpj,omitempty" db:"cnpj"`
	CompanyName       string     `json:"companyName" db:"company_name"`
	CompanyID         int        `json:"companyID,omitempty" db:"company_id"`
	FantasyName       string     `json:"fantasyName" db:"fantasy_name"`
	Phone             string     `json:"phone,omitempty" db:"phone"`
	Email             string     `json:"email,omitempty" db:"email"`
	State             string     `json:"state,omitempty" db:"state"`
	City              string     `json:"city,omitempty" db:"city"`
	District          string     `json:"district,omitempty" db:"district"`
	Street            string     `json:"street,omitempty" db:"street"`
	Number            int        `json:"number,omitempty" db:"number"`
	AddressComplement string     `json:"addressComplement,omitempty" db:"address_complement"`
	Zipcode           string     `json:"zipcode,omitempty" db:"zipcode"`
	Desc              *string    `json:"desc,omitempty" db:"desc"`
	Categories        []Category `json:"categories"`

	Lat      *float64 `db:"lat" json:"lat,omitempty"`
	Lon      *float64 `db:"lon" json:"lon,omitempty"`
	Distance *float64 `db:"distance" json:"distance,omitempty"`

	ActiveApp bool `json:"activeApp" db:"active"`
	Active    bool `db:"active_lyndus" json:"active"`

	PlanValue          decimal.NullDecimal `json:"planValue" db:"plan_value"`
	PlanDay            int                 `json:"planDay" db:"plan_day"`
	BalanceDay         int                 `json:"balanceDay" db:"balance_day"`
	RateAnticipation   decimal.Decimal     `json:"rateAnticipation" db:"rate_anticipation"`
	RateAnticipationBm decimal.Decimal     `json:"rateAnticipationBm" db:"rate_anticipation_bm"`
}

type BasicBS struct {
	ID          int        `json:"id" db:"id"`
	CNPJ        types.CNPJ `json:"cnpj,omitempty" db:"cnpj"`
	CompanyName string     `json:"companyName" db:"company_name"`
	FantasyName string     `json:"fantasyName" db:"fantasy_name"`

	// entity.BS      `db:"bs"`
	// entity.Company `db:"company"`

	//Categories []entity.Category `json:"categories"`
	// WeekDays   []WeekDay         `db:"week_days" json:"weekDays"`
}

type WithdrawalCashSimple struct {
	ID                int             `json:"id"`
	BSID              int             `json:"bsID"`
	CompanyName       string          `json:"companyName"`
	CNPJ              types.CNPJ      `json:"cnpj"`
	ValueReleased     decimal.Decimal `json:"valueReleased"`
	ValueAnticipation decimal.Decimal `json:"valueAnticipation"`
	TotalValue        decimal.Decimal `json:"totalValue"`
	RequestDate       time.Time       `json:"requestDate"`
	PayDate           *time.Time      `json:"payDate"`
	Requested         bool            `json:"requested"`
	Paid              *bool           `json:"paid"`
}

type WithdrawalCash struct {
	ID                int             `json:"id"`
	BSID              int             `json:"bsID"`
	CompanyName       string          `json:"companyName"`
	CNPJ              types.CNPJ      `json:"cnpj"`
	BankDetailsID     *int            `json:"-"`
	Bank              string          `json:"bank"`
	Agency            string          `json:"agency"`
	BankAccount       string          `json:"bankAccount"`
	ChavePIX          string          `json:"chavePIX"`
	CreatedAt         time.Time       `json:"-"`
	ValueReleased     decimal.Decimal `json:"valueReleased"`
	ValueAnticipation decimal.Decimal `json:"valueAnticipation"`
	RateAnticipation  decimal.Decimal `json:"rateAnticipation"`
	TotalValue        decimal.Decimal `json:"totalValue"`
	RequestDate       time.Time       `json:"requestDate"`
	PayDate           *time.Time      `json:"payDate"`
	Requested         bool            `json:"requested"`
	Paid              *bool           `json:"paid"`
	BalanceDay        int             `json:"balanceDay"`
	AccountMovementID *int            `json:"-"`
	AccountName       string          `json:"accountName"`
	AccountReceiptID  int             `json:"-"`
}

type WithdrawalCashBalance struct {
	BalanceReceivableID int             `json:"balanceReceivableID"`
	WithdrawalCashID    int             `json:"withdrawalCashID"`
	Value               decimal.Decimal `json:"value"`
}
