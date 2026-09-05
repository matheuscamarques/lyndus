package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type BSSupplierRepository struct {
	db.Connector
}

func (bsr BSSupplierRepository) GetConnection() *sqlx.DB {
	return bsr.Conn
}

func (bsr BSSupplierRepository) GetConnector() db.Connector {
	return bsr.Connector
}

func (bsr BSSupplierRepository) GetTable() string {
	return bsr.Table
}

func NewBSSupplierRepository() (repo BSSupplierRepository) {
	repo.Table = "bs_supplier"
	repo.UsePostgres()
	return repo
}

//
//func (bsr BSSupplierRepository) SupplierList(name, cnpj string) (supplies []entity.Supplier, err error) {
//	command := `SELECT id,
//					   name,
//					   cnpj
//				FROM bs_supplier
//					WHERE name LIKE $1 OR cnpj LIKE $2`
//	err = bsr.Conn.Select(&supplies, command, name, cnpj)
//	if err == sql.ErrNoRows {
//		err = nil
//	}
//	return supplies, err
//}

func (bsr BSSupplierRepository) GetAll(bsID int, criteria criteria.Criteria)(supplies []aggregate.BSSupplierAggregate, err error){
	command := `SELECT id,
					   name,
					   cnpj
					FROM bs_supplier
					WHERE bs_id=$1
					ORDER BY name`
	command = criteria.Query(command)

	err = bsr.Conn.Select(&supplies, command, bsID)
	return supplies, err
}
func (bsr BSSupplierRepository) GetByID(bsID, supplierID int) (supplier entity.Supplier, err error) {
	command := `SELECT  id,
       				    name,
       					cnpj,
       					phone,
       					representative,
       					obs,
       					address
					FROM bs_supplier
					WHERE bs_id=$1 AND id=$2;`
	err = bsr.Conn.Get(&supplier, command, bsID, supplierID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return supplier, err
}

func (bsr BSSupplierRepository) Create(supplier entity.Supplier) (id int, err error) {
	command := `INSERT INTO bs_supplier(bs_id,
								       name,
								       cnpj,
								       phone,
								       representative,
								       obs,
								       address,
								       created_by) 
            		VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING ID`
	stmt, err := bsr.Conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(supplier.BsID,
		supplier.Name,
		supplier.CNPJ,
		supplier.Phone,
		supplier.Representative,
		supplier.Obs,
		supplier.Address,
		supplier.CreatedBy).Scan(&id)
	return id, err
}

func (bsr BSSupplierRepository) Update(supplier entity.Supplier) error {
	command := `UPDATE bs_supplier SET name=$3,
                       				   cnpj=$4,
                       				   phone=$5,
                       				   representative=$6,
                       				   obs=$7,
                       				   address=$8
						WHERE bs_id=$1 AND id=$2`
	stmt, err := bsr.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(supplier.BsID,
			supplier.ID,
			supplier.Name,
			supplier.CNPJ,
			supplier.Phone,
			supplier.Representative,
			supplier.Obs,
			supplier.Address)
	return err
}