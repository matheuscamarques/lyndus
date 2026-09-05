package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
)

type BSProductRepository struct {
	db.Connector
}

func NewBSProductRepository() (repo BSProductRepository) {
	repo = BSProductRepository{}
	repo.Table = "bs_product"
	repo.UsePostgres()
	return repo
}

func (p BSProductRepository) GetCategories() (categories []entity.ProductCategory, err error) {
	command := `SELECT id,
					   name
					FROM product_category`
	err = p.Conn.Select(&categories, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return categories, err
}

func (p BSProductRepository) GetPackages() (packages []entity.ProductPackage, err error) {
	command := `SELECT id,
					   name
					FROM product_package`
	err = p.Conn.Select(&packages, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return packages, err
}

func (p BSProductRepository) GetUnits() (units []entity.ProductUnit, err error) {
	command := `SELECT id,
					   name
					FROM product_unit`
	err = p.Conn.Select(&units, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return units, err
}

func (p BSProductRepository) GetClasses() (classes []entity.ProductClass, err error) {
	command := `SELECT id,
					   name
					FROM product_class`
	err = p.Conn.Select(&classes, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return classes, err
}

func (p BSProductRepository) GetByID(bsID, productID int) (product entity.Product, err error) {
	command := `SELECT  bsp.id,
       				    bsp.product_class_id,
       					bsp.product_unit_id,
       					bsp.product_package_id,
       					bsp.product_category_id,
       					bsp.product_brand_id,
       					pb.name as brand_name,
       					bsp.bs_supplier_id,
       					bsp.code,
       					bsp.description,
       					bsp.bar_code,
       					bsp.manufacturer_code,
       					bsp.profit_margin,
       					bsp.minimum_stock,
       					bsp.maximum_stock,
       					bsp.quantity,
       					bsp.cost_price,
       					bsp.commission_bm,
       					bsp.sales_price
					FROM bs_product bsp
					INNER JOIN product_brand pb on bsp.product_brand_id = pb.id
					WHERE bsp.bs_id=$1 AND bsp.id=$2;`
	err = p.Conn.Get(&product, command, bsID, productID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return product, err
}

func (p BSProductRepository) CountProducts(bsID int) (n int, err error) {
	command := `SELECT  count(bsp.id)
					FROM bs_product bsp
					WHERE bsp.bs_id=$1`
	err = p.Conn.Get(&n, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return n, err
}

func (p BSProductRepository) Create(product entity.Product) (id int, err error) {
	command := `INSERT INTO bs_product(bs_id,
								       product_class_id,
								       product_unit_id,
								       product_package_id,
								       product_category_id,
								       product_brand_id,
								       bs_supplier_id,
								       code,
								       description,
								       bar_code,
								       manufacturer_code,
								       profit_margin,
								       minimum_stock,
								       maximum_stock,
                       				   quantity,
                       				   cost_price,
                       				   sales_price,
								       commission_bm,
                       				   product_status_id,
									   created_by) 
            		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20) RETURNING id`
	stmt, err := p.Conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(product.BSID,
		product.ProductClassID,
		product.ProductUnitID,
		product.ProductPackageID,
		product.ProductCategoryID,
		product.ProductBrandID,
		product.SupplierID,
		product.Code,
		product.Description,
		product.BarCode,
		product.ManufacturerCode,
		product.ProfitMargin,
		product.MinimumStock,
		product.MaximumStock,
		product.Quantity,
		product.CostPrice,
		product.SalesPrice,
		product.CommissionBM,
		product.Status,
		product.CreatedBy).Scan(&id)
	return id, err
}

func (p BSProductRepository) Update(product entity.Product) error {
	command := `UPDATE bs_product SET product_class_id=$3,
                      				  product_unit_id=$4,
                      				  product_package_id=$5,
                      				  product_category_id=$6,
                      				  product_brand_id=$7,
                      				  bs_supplier_id=$8,
                      				  description=$9,
                      				  bar_code=$10,
                      				  manufacturer_code=$11,
                      				  profit_margin=$12,
                      				  minimum_stock=$13,
                      				  maximum_stock=$14,
                      				  commission_bm=$15,
                      				  quantity=$16,
                      				  cost_price=$17,
                      				  sales_price=$18
						WHERE bs_id=$1 AND id=$2`
	stmt, err := p.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.BSID,
		product.ID,
		product.ProductClassID,
		product.ProductUnitID,
		product.ProductPackageID,
		product.ProductCategoryID,
		product.ProductBrandID,
		product.SupplierID,
		product.Description,
		product.BarCode,
		product.ManufacturerCode,
		product.ProfitMargin,
		product.MinimumStock,
		product.MaximumStock,
		product.CommissionBM,
		product.Quantity,
		product.CostPrice,
		product.SalesPrice)
	return err
}

func (p BSProductRepository) UpdateProductStock(product entity.Product) error {
	command := `UPDATE bs_product SET quantity=$3,
                      				  cost_price=$4,
                      				  profit_margin=$5,
                      				  sales_price=$6
						WHERE bs_id=$1 AND id=$2`
	stmt, err := p.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.BSID,
		product.ID,
		product.Quantity,
		product.CostPrice,
		product.ProfitMargin,
		product.SalesPrice)
	return err
}

func (p BSProductRepository) GetAllProducts(bsID int, c criteria.Criteria) (products []entity.Product, err error) {
	command := `SELECT bsp.id,
					   bsp.code,
       				   pb.name as brand_name,
					   bsp.description,
					   bsp.quantity,
					   bsp.cost_price,
					   bsp.sales_price
					FROM bs_product bsp
					INNER JOIN product_brand pb on bsp.product_brand_id = pb.id
					WHERE bs_id = $1`
	command = c.Query(command)
	err = p.Conn.Select(&products, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return products, err
}
func (p BSProductRepository) GetProductSimple(bsID, productID int) (product entity.ProductSimple, err error) {
	command := `SELECT bsp.id,
       				   pb.name as brand_name,
					   bsp.description,
					   bsp.cost_price,
       				   bsp.profit_margin,
       				   bsp.sales_price
					FROM bs_product bsp
					INNER JOIN product_brand pb on bsp.product_brand_id = pb.id
					WHERE bs_id = $1 AND bsp.id=$2`
	err = p.Conn.Get(&product, command, bsID, productID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return product, err
}

func (p BSProductRepository) GetProductsSimple(bsID int) (products []entity.ProductSimple, err error) {
	command := `SELECT bsp.id,
       				   bsp.product_class_id,
       				   pb.name as brand_name,
					   bsp.description,
					   bsp.cost_price,
       				   bsp.profit_margin,
       				   bsp.sales_price
					FROM bs_product bsp
					INNER JOIN product_brand pb on bsp.product_brand_id = pb.id
					WHERE bs_id = $1`
	err = p.Conn.Select(&products, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return products, err
}

func (p BSProductRepository) StockSave(bsId, createdBy, id int, quantity int, costPrice, salesPrice, profitMargin decimal.Decimal) error {
	product, err := p.GetByID(bsId, id)
	if err != nil {
		return err
	}
	if product.ID == 0 {
		return p.Error(errors.New("bs_product bs not founded id == 0"))
	}
	command := `UPDATE bs_product SET quantity = quantity + $3,
	                			          cost_price = $4,
										  sales_price = $5,
                      					  profit_margin = $6
									WHERE bs_id = $1 AND id = $2`
	stmt, err := p.Conn.Prepare(command)
	if err != nil {
		return p.Error(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsId, id, quantity, costPrice, salesPrice, profitMargin)
	if err != nil {
		return p.Error(err)
	}

	command = `INSERT INTO bs_product_stock_history(
                                     		 bs_product_id,
                                             new_cost_price,
                                             new_quantity,
                                     		 new_sales_price,
                                             old_cost_price,
                                             old_quantity,
                                     		 old_sales_price,
                                             created_by
                                             )
                                             VALUES($1,$2,$3,$4,$5,$6,$7,$8)
                                             RETURNING id
                                             `
	stmt, err = p.Conn.Prepare(command)
	defer stmt.Close()
	var rID int
	err = stmt.QueryRow(
		product.ID,
		costPrice,
		product.Quantity+quantity,
		salesPrice,
		product.CostPrice,
		product.Quantity,
		product.SalesPrice,
		createdBy,
	).Scan(&rID)
	return err
}
