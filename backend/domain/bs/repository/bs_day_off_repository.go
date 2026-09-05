package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type BSDayOffRepository struct {
	contracts.BSDayOffRepositoryInterface
	conn *sqlx.DB
}

func NewBSDayOffRepository() BSDayOffRepository {
	return BSDayOffRepository{
		conn: postgres.DB,
	}
}

func (do BSDayOffRepository) AddBSDayOff(bsID int, day entity.DayOff) (id int, err error) {
	command := `INSERT INTO bs_day_off(bs_id, start_date, end_date ) VALUES($1, $2, $3) RETURNING ID`
	stmt, err := do.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(bsID, day.StartDate, day.EndDate).Scan(&id)
	return id, err
}

func (do BSDayOffRepository) GetBSDaysOff(bsID int) (days []entity.DayOff, err error) {
	command := `SELECT id, 
       				   start_date,
       				   end_date
                FROM bs_day_off
                WHERE bs_id=$1`
	err = do.conn.Select(&days, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return days, err
}

func (do BSDayOffRepository) DeleteBSDayOff(bsID, dayOffID int) error {
	command := `DELETE FROM bs_day_off WHERE bs_id = $1 AND id = $2`
	stmt, err := do.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, dayOffID)
	return err
}

func (do BSDayOffRepository) AddBSBMDayOff(bsID, bmID int, day entity.DayOff) (id int, err error) {
	command := `INSERT INTO bs_bm_day_off(bs_id, bs_bm_id, start_date, end_date ) VALUES($1, $2, $3, $4) RETURNING ID`
	stmt, err := do.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(bsID, bmID, day.StartDate, day.EndDate).Scan(&id)
	return id, err
}

func (do BSDayOffRepository) GetBSBMDaysOff(bsID, bmID int) (days []entity.DayOff, err error) {
	command := `SELECT id, 
       				   start_date,
       				   end_date
                FROM bs_bm_day_off
                WHERE bs_id=$1 AND bs_bm_id=$2`
	err = do.conn.Select(&days, command, bsID, bmID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return days, err
}

func (do BSDayOffRepository) DeleteBSBMDayOff(bsID, bmID, dayOffID int) error {
	command := `DELETE FROM bs_bm_day_off WHERE bs_id = $1 AND bs_bm_id=$2 AND id = $3`
	stmt, err := do.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, bmID, dayOffID)
	return err
}
