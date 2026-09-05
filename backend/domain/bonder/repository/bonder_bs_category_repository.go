package repository

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/infra/db"
)

type BonderBSCategoryRepository struct {
	db.Connector
}

func NewBonderBSCategoryRepository() (repo BonderBSCategoryRepository) {
	repo.Table = "bs_category"
	repo.UsePostgres()
	return repo
}

func (bcr BonderBSCategoryRepository) GetAll() (categories []entity.Category, err error) {
	command := `SELECT id,
					   "desc"
					FROM bs_category`
	err = bcr.Conn.Select(&categories, command)
	return categories, err
}

func (bcr BonderBSCategoryRepository) GetByBS(bsID int) (categories []entity.Category, err error) {
	command := `SELECT bsc.id,
       				   bsc."desc"
					FROM bs_bs_category bbc
					INNER JOIN bs_category bsc on bbc.bs_category_id = bsc.id
					WHERE bbc.bs_id=$1`
	err = bcr.Conn.Select(&categories, command, bsID)

	return categories, err
}

func (bcr BonderBSCategoryRepository) Create(bsID, bsCategoryID int) error {
	command := `INSERT INTO bs_bs_category(bs_id, bs_category_id) VALUES($1,$2) `
	stmt, err := bcr.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(bsID, bsCategoryID)
	return err
}

func (bcr BonderBSCategoryRepository) Remove(bsID, bsCategoryID int) error {
	command := `DELETE FROM bs_bs_category WHERE bs_id = $1 AND bs_category_id = $2`
	_, err := bcr.Conn.Exec(command, bsID, bsCategoryID)
	return err
}
