package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type BSProductBrandRepository struct {
	db.Connector
}

func (bsr BSProductBrandRepository) GetConnection() *sqlx.DB {
	return bsr.Conn
}

func (bsr BSProductBrandRepository) GetConnector() db.Connector {
	return bsr.Connector
}

func (bsr BSProductBrandRepository) GetTable() string {
	return bsr.Table
}

func NewBSProductBrandRepository() (repo BSProductBrandRepository) {
	repo.Table = "product_brand"
	repo.UsePostgres()
	return repo
}

func (bsr BSProductBrandRepository) GetAll(criteria criteria.Criteria) (brands []aggregate.BSProductBrandAggregate, err error) {
	command := criteria.Query(`SELECT id,
					   						name
									FROM product_brand
									ORDER BY name`)

	err = bsr.Conn.Select(&brands, command)
	return brands, err
}

func (bsr BSProductBrandRepository) GetByName(name string) (brand entity.ProductBrand, err error) {
	command := `SELECT id,
					   name
				FROM product_brand
				WHERE name = $1`

	err = bsr.Conn.Get(&brand, command, name)
	if err == sql.ErrNoRows {
		err = nil
	}
	return brand, err
}

func (bsr BSProductBrandRepository) GetByID(brandID int) (brand entity.ProductBrand, err error) {
	command := `SELECT  id,
       				    name
					FROM product_brand
					WHERE id=$1`
	err = bsr.Conn.Get(&brand, command, brandID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return brand, err
}

func (bsr BSProductBrandRepository) SearchByName(name string) (brands []entity.ProductBrand, err error) {
	command := `SELECT  id,
       				    name
					FROM product_brand
					WHERE LOWER(name) LIKE LOWER($1) ORDER BY name`
	err = bsr.Conn.Select(&brands, command, name+"%")
	if err == sql.ErrNoRows {
		err = nil
	}
	return brands, err
}

func (bsr BSProductBrandRepository) Create(product entity.ProductBrand) (id int, err error) {
	command := `INSERT INTO product_brand(name) 
            			VALUES($1) RETURNING id`
	stmt, err := bsr.Conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(product.Name).Scan(&id)
	return id, err
}

func (bsr BSProductBrandRepository) Update(brand entity.ProductBrand) error {
	command := `UPDATE bs_supplier SET name=$2 WHERE id=$1`
	stmt, err := bsr.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(brand.ID, brand.Name)
	return err
}
