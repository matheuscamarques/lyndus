package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type BSWeekDayRepository struct {
	contracts.BSWeekDayRepositoryInterface
	conn *sqlx.DB
}

func NewBSWeekDayRepository() BSWeekDayRepository {
	return BSWeekDayRepository{
		conn: postgres.DB,
	}
}

func (w BSWeekDayRepository) GetWeekDays() (weekDays []entity.WeekDay, err error) {
	command := `SELECT id, nome as name FROM weekday ORDER BY sequence`
	err = w.conn.Select(&weekDays, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (w BSWeekDayRepository) Get(bsID, weekDayID int) (wds []entity.WeekDayBS, err error) {
	command := `SELECT id,
                  bs_id,
                  week_day_id,
                  start_time,
                  end_time
            FROM bs_week_day
            WHERE bs_id=$1 AND week_day_id=$2 ORDER BY start_time`
	err = w.conn.Select(&wds, command, bsID, weekDayID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return wds, err
}
