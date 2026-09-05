package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type AppUserCategoryRepository struct {
	contracts.AppUserCategoryRepositoryInterface
	conn *sqlx.DB
}

func NewAppUserCategoryRepository() *AppUserCategoryRepository {
	return &AppUserCategoryRepository{
		conn: postgres.DB,
	}
}

func (aucr AppUserCategoryRepository) GetCategory(c *entity.Category) error {
	command := `SELECT bca.id, bca.desc FROM bs_category bca`
	err := aucr.conn.Get(c, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return err
}

func (aucr AppUserCategoryRepository) GetCategories() (categories []entity.Category, err error) {
	command := `SELECT bca.id, bca.desc FROM bs_category bca `
	err = aucr.conn.Select(&categories, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return categories, err
}
