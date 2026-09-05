package repository

import (
	"bitbucket.org/lyndus/backend/global/aggregate"
	"database/sql"
	"time"

	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/constants"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type AppUserBMRepository struct {
	conn *sqlx.DB
}

func NewAppUserBMRepository() *AppUserBMRepository {
	return &AppUserBMRepository{
		conn: postgres.DB,
	}
}

//
//func (aus AppUserServiceRepository) GetServiceCategory() (serviceCategories []entity.ServiceCategory, err error) {
//	command := `SELECT id, "desc" FROM bs_service_category ORDER BY id `
//	err = aus.conn.Select(&serviceCategories, command)
//	if err == sql.ErrNoRows {
//		err = nil
//	}
//	return serviceCategories, err
//}

func (aub AppUserBMRepository) GetBMByID(bsID, bmID int) (bm entity.BM, err error) {
	command := `SELECT bbm.id,
       				   bbm."name",
       				   bbm.obs,
       				   bbm.service_step_time
					FROM  bs_bm bbm 
					WHERE bbm.bs_id=$1 AND bbm.id=$2 AND bbm.bm_status_id=$3
					ORDER BY name `
	err = aub.conn.Get(&bm, command, bsID, bmID, constants.BMStatusActive)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bm, err
}

func (aub AppUserBMRepository) GetBMByIDAndDuration(bsID, bmID, serviceID int) (bm entity.BM, err error) {
	command := `SELECT bbm.id,
       				   bbm."name",
       				   bbm.obs,
       				   bbm.service_step_time,
					   bbs.duration
					FROM  bs_bm bbm 
					INNER JOIN bs_bm_service bbs on bbm.id = bbs.bs_bm_id
					WHERE bbm.bs_id=$1 
					  AND bbm.id=$2 
					  AND bbm.bm_status_id=$3
					  AND bbs.bs_service_id=$4
					ORDER BY name `
	err = aub.conn.Get(&bm, command, bsID, bmID, constants.BMStatusActive, serviceID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bm, err
}

func (aub AppUserBMRepository) GetBMsFullByWeekDay(bsID, serviceID, weekDayID int) (bms []entity.BM, err error) {
	command := `SELECT bbm.id,
       				   bbm."name",
       				   bbm.obs,
       				   bbm.service_step_time,
					   bbs.duration,
       				   bwd.start_time, 
       				   bwd.end_time
					FROM  bs_bm bbm 
					INNER JOIN bs_bm_service bbs on bbm.id = bbs.bs_bm_id
					INNER JOIN bs_bm_week_day bwd on bbm.id = bwd.bs_bm_id
					WHERE bbm.bs_id=$1 
					  AND bbm.bm_status_id=$2
					  AND bbs.bs_service_id=$3
					  AND bwd.week_day_id = $4
					ORDER BY name `
	err = aub.conn.Select(&bms, command, bsID, constants.BMStatusActive, serviceID, weekDayID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bms, err
}

func (aub AppUserBMRepository) GetBMFullByWeekDay(bsID, serviceID, bmID, weekDayID int) (bm entity.BM, err error) {
	command := `SELECT bbm.id,
       				   bbm."name",
       				   bbm.obs,
       				   bbm.service_step_time,
					   bbs.duration,
       				   bwd.start_time, 
       				   bwd.end_time
					FROM  bs_bm bbm 
					INNER JOIN bs_bm_service bbs on bbm.id = bbs.bs_bm_id
					INNER JOIN bs_bm_week_day bwd on bbm.id = bwd.bs_bm_id
					WHERE bbm.bs_id=$1 
					  AND bbm.id=$2
					  AND bbm.bm_status_id=$3
					  AND bbs.bs_service_id=$4
					  AND bwd.week_day_id = $5
					ORDER BY name `
	err = aub.conn.Get(&bm, command, bsID, bmID, constants.BMStatusActive, serviceID, weekDayID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bm, err
}

func (aub AppUserBMRepository) GetBMsBasicByServiceID(bsID, serviceID int) (bms []entity.BM, err error) {
	command := `SELECT bbm.id,   		
       				   bbm."name",
       				   bbm.obs,
					   bbm.desc,
       				   bbs.duration
					FROM bs_bm_service bbs
					INNER JOIN bs_bm bbm ON bbs.bs_bm_id = bbm.id
					INNER JOIN bs_bm_week_day bwd on bbm.id = bwd.bs_bm_id
					WHERE bbm.bs_id=$1 
					  AND bs_service_id=$2 
					  AND bbm.bm_status_id=$3
					GROUP BY BBM.id, bbm."name", bbm.obs, bbm.desc, bbs.duration
					ORDER BY name`
	err = aub.conn.Select(&bms, command, bsID, serviceID, constants.BMStatusActive)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bms, err
}

func (aub AppUserBMRepository) GetBMsByServiceID(bsID, serviceID int) (bms []entity.BM, err error) {
	command := `SELECT bbm.id,   		
       				   bbm."name",
       				   bbm.obs,
					   bbm.service_step_time,
					   bbm.desc,
       				   bbs.duration
					FROM bs_bm_service bbs
					INNER JOIN bs_bm bbm ON bbs.bs_bm_id = bbm.id
					WHERE bbm.bs_id=$1 
					  AND bs_service_id=$2 
					  AND bbm.bm_status_id=$3
					ORDER BY name`
	err = aub.conn.Select(&bms, command, bsID, serviceID, constants.BMStatusActive)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bms, err
}

func (aub AppUserBMRepository) GetBMWeekDays(bmID int) (weekDays []entity.WeekDayBS, err error) {
	command := `SELECT wed.id,
					   bwd.start_time,
					   bwd.end_time,
       				   bwd.week_day_id
				FROM weekday wed
				INNER JOIN bs_bm_week_day bwd ON wed.id = bwd.week_day_id 
				WHERE bwd.bs_bm_id=$1 
				ORDER BY wed.sequence`
	err = aub.conn.Select(&weekDays, command, bmID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return weekDays, err
}

func (aub AppUserBMRepository) GetBMWeekDaysWithName(bmID int) (weekDays []aggregate.WeekDay, err error) {
	command := `SELECT wed.id,
					   bwd.start_time,
					   bwd.end_time,
       				   wed.nome as name
				FROM weekday wed
				INNER JOIN bs_bm_week_day bwd ON wed.id = bwd.week_day_id
				WHERE bwd.bs_bm_id=$1 ORDER BY wed.sequence`
	err = aub.conn.Select(&weekDays, command, bmID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return weekDays, err
}

func (aub AppUserBMRepository) GetBMDaysOff(bsID, bmID int, start, end time.Time) (days []entity.DayOff, err error) {
	command := `SELECT id, 
       				   start_date,
       				   end_date
                FROM bs_bm_day_off
                WHERE bs_id=$1 
                	AND bs_bm_id=$2
                	AND start_date >= $3
  					AND end_date <= $4
                  	OR (DATE($3) between DATE(start_date) AND DATE(end_date)
                           OR DATE($4) BETWEEN DATE(start_date) AND DATE(end_date))`
	err = aub.conn.Select(&days, command, bsID, bmID, start, end)
	if err == sql.ErrNoRows {
		err = nil
	}
	return days, err
}
