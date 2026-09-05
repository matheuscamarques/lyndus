package repository

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type AppUserBalanceRepository struct {
	conn *sqlx.DB
}

func NewAppUserBalanceRepository() AppUserBalanceRepository {
	return AppUserBalanceRepository{
		conn: postgres.DB,
	}
}

func (w AppUserBalanceRepository) GetBalance(appUserID int) (balance entity.Balance, err error) {
	command := `SELECT id, value FROM app_user_benefit WHERE app_user_id=$1`
	err = w.conn.Get(&balance, command, appUserID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (w AppUserBalanceRepository) UpdateBalance(appUserID int, value decimal.Decimal) error {
	command := `UPDATE app_user_benefit SET value=$2, last_update=now() WHERE app_user_id=$1`
	stmt, err := w.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(appUserID, value)
	return err
}

func (w AppUserBalanceRepository) InsertBalanceHistory(appUserID int, oldValue, newValue decimal.Decimal) (int, error) {
	command := "INSERT INTO public.app_user_benefit_history(app_user_id, value_old, value_new) VALUES($1,$2,$3) RETURNING id"
	stmt, err := w.conn.Prepare(command)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(appUserID, oldValue, newValue).Scan(&id)
	return id, err
}
