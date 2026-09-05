package repository

import (
	"bitbucket.org/lyndus/backend/global/basic"
	"database/sql"
	"time"

	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/constants"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/jmoiron/sqlx"
)

type AppUserScheduleRepository struct {
	contracts.AppUserScheduleRepositoryInterface
	conn *sqlx.DB
}

func NewAppUserScheduleRepository() *AppUserScheduleRepository {
	return &AppUserScheduleRepository{
		conn: postgres.DB,
	}
}

func (s AppUserScheduleRepository) ChangeStatusSchedule(bsID, scheduleID, status int) (err error) {
	command := `UPDATE bs_schedule_time SET schedule_status_id =$3
					WHERE bs_id = $1 AND id = $2`
	stmt, err := s.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, scheduleID, status)
	return err
}

func (s AppUserScheduleRepository) GetSchedulingsByAppUserID(appUserID, activePage, itemsPerPage int, date types.Date) (scheduling []basic.AppScheduling, totalPages int, err error) {
	command := `SELECT COUNT(*)
				FROM bs_schedule_time bst
						 INNER JOIN bs_person bpe on bpe.id = bst.bs_person_id
						 INNER JOIN app_user aus on aus.id = bpe.app_user_id
						 WHERE aus.id=$1`

	if !date.IsZero() {
		command = command + ` AND DATE(bst.start_time)=DATE($2) `
		err = s.conn.Get(&totalPages, command, appUserID, date)
	} else {
		err = s.conn.Get(&totalPages, command, appUserID)
	}
	if err != nil {
		return scheduling, totalPages, err
	}

	rest := totalPages % itemsPerPage
	totalPages = totalPages / itemsPerPage

	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT bst.id,
					   bst.bs_id,
					   bst.bs_bm_id,
					   bbm.name as bs_bm_name,
					   bst.bs_service_id,
					   bse.name as bs_service_name,
					   bst.schedule_status_id,
					   bst.start_time,
					   bst.end_time,
					   bst.created_at
				FROM bs_schedule_time bst
						 INNER JOIN bs_bm bbm on bst.bs_bm_id = bbm.id
						 INNER JOIN bs_service bse on bst.bs_service_id = bse.id
						 INNER JOIN bs_person bpe on bpe.id = bst.bs_person_id
						 INNER JOIN app_user aus on aus.id = bpe.app_user_id
						 WHERE aus.id=$1 `

	if !date.IsZero() {
		command = command + ` AND DATE(bst.start_time)=DATE($2) 
						 		ORDER By bst.start_time DESC
						  		OFFSET $3 LIMIT $4`
		err = s.conn.Select(&scheduling, command, appUserID, date, (activePage-1)*itemsPerPage, itemsPerPage)
	} else {
		command = command + ` ORDER By bst.start_time DESC
						  		OFFSET $2 LIMIT $3`
		err = s.conn.Select(&scheduling, command, appUserID, (activePage-1)*itemsPerPage, itemsPerPage)
	}
	if err == sql.ErrNoRows {
		err = nil
	}
	return scheduling, totalPages, err
}

func (s AppUserScheduleRepository) GetSchedulingAppUser(appUserID, scheduleID int) (scheduling entity.Scheduling, err error) {
	command := `SELECT bst.id,
       				   bst.bs_id,
       				   bst.bs_bm_id,
       				   bst.bs_person_id,
       				   bpe.name as bs_person_name,
       				   bst.bs_service_id,
       				   bst.schedule_status_id,
       				   bst.start_time,
       				   bst.end_time,
					   bst.app_user,
       				   bst.created_at
       				FROM bs_schedule_time bst
						 INNER JOIN bs_person bpe on bpe.id = bst.bs_person_id
						 INNER JOIN app_user aus on aus.id = bpe.app_user_id
						 WHERE aus.id=$1
						 AND bst.id=$2`
	err = s.conn.Get(&scheduling, command, appUserID, scheduleID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return scheduling, err
}

func (s AppUserScheduleRepository) GetBSScheduling(bsId int, date time.Time) (scheduling []entity.Scheduling, err error) {
	command := `SELECT bst.id,
       				   bst.bs_id,
       				   bst.bs_bm_id,
       				   bbm.name as bs_bm_name,
       				   bst.bs_person_id,
       				   bpe.name as bs_person_name,
       				   bst.bs_service_id,
       				   bse.name as bs_service_name,
       				   bst.schedule_status_id,
       				   bst.start_time,
       				   bst.end_time,
					   bst.app_user,
       				   bst.created_at
       				FROM bs_schedule_time bst
					INNER JOIN bs_bm bbm on bst.bs_bm_id = bbm.id
					INNER JOIN bs_service bse on bst.bs_service_id = bse.id
					INNER JOIN bs_person bpe on bpe.id = bst.bs_person_id
					WHERE bst.bs_id = $1
					  AND bst.start_time >= $2
  					  AND bst.start_time < $3
       				ORDER BY bst.start_time`
	err = s.conn.Select(&scheduling, command, bsId, date.Truncate(time.Hour*24), date.Add(time.Hour*24))
	if err == sql.ErrNoRows {
		err = nil
	}
	return scheduling, err
}

func (s AppUserScheduleRepository) GetBSBMScheduling(bsId, bmID int, start, end time.Time) (scheduling []entity.SchedulingBAsic, err error) {
	command := `SELECT bst.id,
       				   bst.bs_id,
       				   bst.bs_bm_id,
       				   bst.bs_person_id,
       				   bst.bs_service_id,
       				   bst.schedule_status_id,
       				   bst.start_time,
       				   bst.end_time,
					   bst.app_user,
       				   bst.created_at
       				FROM bs_schedule_time bst
					WHERE bst.bs_id = $1
					  AND bst.bs_bm_id = $2
					  AND bst.schedule_status_id = $3
					  AND DATE(bst.start_time) >= DATE($4)
  					  AND DATE(bst.start_time) <= DATE($5)
       				ORDER BY bst.start_time`
	err = s.conn.Select(&scheduling, command, bsId, bmID, constants.SCHEDULESCHEDULED, start, end)
	if err == sql.ErrNoRows {
		err = nil
	}
	return scheduling, err
}

func (s AppUserScheduleRepository) CreateScheduling(scheduling entity.Scheduling) (id int, err error) {
	command := `INSERT INTO bs_schedule_time (
						bs_id,
						bs_bm_id,
                        bs_person_id,
						bs_service_id,
						schedule_status_id,
						app_user,
						start_time,
						end_time
						)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                       	RETURNING ID`
	stmt, err := s.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(scheduling.BSID,
		scheduling.BSBmID,
		scheduling.BSPersonID,
		scheduling.BSServiceID,
		constants.SCHEDULESCHEDULED,
		scheduling.AppUser,
		scheduling.StartTime,
		scheduling.EndTime,
	).Scan(&id)
	return id, err
}

func (s AppUserScheduleRepository) GetBSScheduledTime(bsId, bmID int, startTime, endTime types.DateTime) (ids []int, err error) {
	command := `SELECT id FROM bs_schedule_time
						WHERE bs_id=$1
				  		AND bs_bm_id=$2
				  		AND schedule_status_id = $3
				  		AND(
								$4 BETWEEN start_time AND end_time
					  		OR
								$5 BETWEEN start_time AND end_time
						)`
	err = s.conn.Select(&ids, command, bsId, bmID, constants.SCHEDULESCHEDULED, startTime, endTime)
	if err == sql.ErrNoRows {
		err = nil
	}
	return ids, err
}
