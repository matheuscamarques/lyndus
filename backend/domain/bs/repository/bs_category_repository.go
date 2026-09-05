package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type BSCategoryRepository struct {
	contracts.BSCategoryRepositoryInterface
	conn *sqlx.DB
}

func NewBSCategoryRepository() BSCategoryRepository {
	return BSCategoryRepository{
		conn: postgres.DB,
	}
}

func (c BSCategoryRepository) GetBSCategories(bsID int) (categories []entity.Category, err error) {
	command := `SELECT bca.id,
                  bca.desc
                FROM bs_category bca 
                INNER JOIN bs_bs_category bbc ON bca.id = bbc.bs_category_id
                WHERE bbc.bs_id=$1`
	err = c.conn.Select(&categories, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return categories, err
}

func (c BSCategoryRepository) AddBSCategory(bsID, bsCategoryID int) error {
	command := `INSERT INTO bs_bs_category(bs_id, bs_category_id) VALUES($1, $2)`
	stmt, err := c.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, bsCategoryID)
	return err
}

func (c BSCategoryRepository) DeleteBSCategory(bsID, bsCategoryID int) error {
	command := `DELETE FROM bs_bs_category WHERE  bs_id= $1 AND bs_category_id = $2`
	stmt, err := c.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, bsCategoryID)
	return err
}
