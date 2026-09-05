package repository

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type AppUserWeekDayRepository struct {
	conn *sqlx.DB
}

func NewAppUserWeekDayRepository() AppUserWeekDayRepository {
	return AppUserWeekDayRepository{
		conn: postgres.DB,
	}
}

func (w AppUserWeekDayRepository) GetWeekDays() (weekDays []entity.WeekDay, err error) {
	command := `SELECT id, nome as name FROM weekday ORDER BY sequence`
	err = w.conn.Select(&weekDays, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (w AppUserWeekDayRepository) Get(bsID, weekDayID int) (wds []entity.WeekDayBS, err error) {
	command := `SELECT id,
                  bs_id,
                  week_day_id,
                  start_time,
                  end_time
            FROM bs_week_day
            WHERE bs_id=$1 AND week_day_id=$2`
	err = w.conn.Select(&wds, command, bsID, weekDayID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return wds, err
}
