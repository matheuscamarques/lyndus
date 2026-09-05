package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
)

type AppUserGenderRepository struct {
	conn *sqlx.DB
}

func NewAppUserGenderRepository() *AppUserGenderRepository {
	return &AppUserGenderRepository{
		conn: postgres.DB,
	}
}

func (augr AppUserGenderRepository) GetGenders() (genders []entity.Gender, err error) {
	command := `SELECT id, "desc" FROM gender `
	err = augr.conn.Select(&genders, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return genders, err
}

func (augr AppUserGenderRepository) GetCategories() (categories []entity.Category, err error) {
	command := `SELECT bca.id, bca.desc FROM bs_category bca `
	err = augr.conn.Select(&categories, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return categories, err
}

func (augr AppUserGenderRepository) GetGenderID(genderID int) (id int, err error) {
	command := `SELECT id FROM gender WHERE id=$1`
	err = augr.conn.Get(&id, command, genderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return id, err
}

func (augr AppUserGenderRepository) GetGenderByID(genderID int) (gender entity.Gender, err error) {
	command := `SELECT id, "desc" FROM gender WHERE id=$1`
	err = augr.conn.Get(&gender, command, genderID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return gender, err
}
