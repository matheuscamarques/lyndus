package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/constants"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/types"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type BSScheduleRepository struct {
	contracts.BSScheduleRepositoryInterface
	conn *sqlx.DB
}

func NewBSScheduleRepository() *BSScheduleRepository {
	return &BSScheduleRepository{
		conn: postgres.DB,
	}
}

func (s BSScheduleRepository) ChangeStatusSchedule(bsID, scheduleID, status int) (err error) {
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

func (s BSScheduleRepository) SetOrderIDSchedule(bsID, scheduleID, orderID int) (err error) {
	command := `UPDATE bs_schedule_time SET bs_order_id =$3
					WHERE bs_id = $1 AND id = $2`
	stmt, err := s.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, scheduleID, orderID)
	return err
}

func (s BSScheduleRepository) GetBSSchedule(bsId, scheduleID int) (scheduling entity.Scheduling, err error) {
	command := `SELECT bst.id,
      				   bst.bs_id,
      				   bst.bs_bm_id,
      				   bst.bs_person_id,     				   
      				   bst.bs_service_id,
      				   bst.schedule_status_id,
      				   bst.start_time,
      				   bst.end_time,
      				   bst.app_user,
      				   bst.created_at,
       				   bst.bs_order_id
      				FROM bs_schedule_time bst
					WHERE bst.bs_id = $1 
					  AND bst.id = $2`
	err = s.conn.Get(&scheduling, command, bsId, scheduleID )
	if err == sql.ErrNoRows {
		err = nil
	}
	return scheduling, err
}





func (s BSScheduleRepository) GetBSScheduling(bsId int, date time.Time) (scheduling []entity.Scheduling, err error) {
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
      				   bst.created_at,
       				   bst.bs_order_id
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



func (s BSScheduleRepository) CreateScheduling(scheduling entity.Scheduling) (id int, err error) {
	command := `INSERT INTO bs_schedule_time (
						bs_id,
						bs_bm_id,
                       bs_person_id,
						bs_service_id,
						schedule_status_id,
						start_time,
						end_time,
						app_user)
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
		scheduling.StartTime,
		scheduling.EndTime,
		scheduling.AppUser).Scan(&id)
	return id, err
}

func (s BSScheduleRepository) GetBSScheduledTime(bsId, bmID int, startTime, endTime types.DateTime) (ids []int, err error) {
	command := `SELECT id FROM bs_schedule_time
						WHERE bs_id=$1
				  		AND bs_bm_id=$2
				  		AND schedule_status_id = $3
				  		AND(
								($4 >= start_time AND $4 <  end_time)
					  		OR
								$5 BETWEEN start_time AND end_time
						);`
	err = s.conn.Select(&ids, command, bsId, bmID, constants.SCHEDULESCHEDULED, startTime, endTime)
	if err == sql.ErrNoRows {
		err = nil
	}
	return ids, err
}
