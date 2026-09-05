package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Product struct {
	ID                int             `json:"id" db:"id"`
	BSID              int             `json:"-" db:"bs_id"`
	ProductClassID    int             `json:"classID" db:"product_class_id"`
	ProductUnitID     int             `json:"unitID" db:"product_unit_id"`
	ProductPackageID  int             `json:"packageID" db:"product_package_id"`
	ProductCategoryID int             `json:"categoryID" db:"product_category_id"`
	ProductBrandID    int             `json:"brandID" db:"product_brand_id"`
	BrandName         string          `json:"brandName" db:"brand_name"`
	SupplierID        int             `json:"supplierID" db:"bs_supplier_id"`
	Code              string          `json:"code" db:"code"`
	Description       string          `json:"description" db:"description"`
	BarCode           string          `json:"barCode" db:"bar_code"`
	ManufacturerCode  string          `json:"manufacturerCode" db:"manufacturer_code"`
	ProfitMargin      decimal.Decimal `json:"profitMargin" db:"profit_margin"`
	MinimumStock      int             `json:"minimumStock" db:"minimum_stock"`
	MaximumStock      int             `json:"maximumStock" db:"maximum_stock"`
	CommissionBM      bool            `json:"commissionBM" db:"commission_bm"`
	Quantity          int             `json:"quantity" db:"quantity"`
	CostPrice         decimal.Decimal `json:"costPrice" db:"cost_price"`
	SalesPrice        decimal.Decimal `json:"salesPrice" db:"sales_price"`
	CreatedBy         int             `json:"createdBy,omitempty" db:"created_by"`
	CreatedAt         time.Time       `json:"createdAt,omitempty" db:"created_at"`
	Status            int             `json:"status" db:"product_status_id"`
}

type ProductSimple struct {
	ID             int             `json:"id" db:"id"`
	ProductClassID int             `json:"product_class_id" db:"product_class_id"`
	Description    string          `json:"description" db:"description"`
	BrandName      string          `json:"brandName" db:"brand_name"`
	Value          decimal.Decimal `json:"value" db:"value"`

	SalesPrice   decimal.Decimal `json:"salesPrice" db:"sales_price"`
	ProfitMargin decimal.Decimal `json:"-" db:"profit_margin"`
	CostPrice    decimal.Decimal `json:"-" db:"cost_price"`
}
