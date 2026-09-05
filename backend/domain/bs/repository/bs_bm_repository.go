package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db"
	"bitbucket.org/lyndus/backend/infra/types"
	"database/sql"
)

type BSBMRepository struct {
	db.Connector
}

func NewBSBMRepository() (repo BSBMRepository) {
	repo.Table = "bs_bm"
	repo.UsePostgres()
	return repo
}

func (b BSBMRepository) GetBMService(bmID, serviceID int) (service entity.Service, err error) {
	//Todo VERIFICAR SQL
	command := `SELECT bs_bm_id,
       				   bs_service_id as id,
       				   duration
					FROM bs_bm_service
					WHERE bs_bm_id=$1 
					  AND bs_service_id =$2 `
	err = b.Conn.Get(&service, command, bmID, serviceID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return service, err
}

func (b BSBMRepository) GetBmID(bsID, bmID, bmStatusID int) (n int, err error) {
	command := `SELECT bs_bm.id
					FROM bs_bm bs_bm
					WHERE bs_bm.bs_id=$1 
					  AND bs_bm.id=$2 
					  AND bm_status_id=$3`
	err = b.Conn.Get(&n, command, bsID, bmID, bmStatusID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return n, err
}

func (b BSBMRepository) GetBM(bsID, bmID, bmStatusID int) (bm entity.BM, err error) {

	command := `SELECT bs_bm.id,
					   bs_bm.name,
					   bs_bm.cpf,
       				   bs_bm.cnpj,
					   bs_bm.obs,
					   bs_bm.commission_service,
					   bs_bm.commission_product,
					   bs_bm.commission_supply,
					   bs_bm.service_step_time,
       				   bs_bm.desc
					FROM bs_bm bs_bm
					WHERE bs_bm.bs_id=$1 
					  AND bs_bm.id=$2 
					  AND bm_status_id=$3`
	err = b.Conn.Get(&bm, command, bsID, bmID, bmStatusID)
	if err == sql.ErrNoRows {
		return bm, nil
	}
	return bm, err
}

func (b BSBMRepository) GetBmIDByCPF(bsID int, document types.CPF) (id int, err error) {
	command := `SELECT bs_bm.id 
					FROM bs_bm bs_bm 
					WHERE bs_bm.bs_id = $1 
					  AND bs_bm.cpf=$2`
	err = b.Conn.Get(&id, command, bsID, document)
	if err == sql.ErrNoRows {
		err = nil
	}
	return id, err
}

func (b BSBMRepository) CreateBM(bm entity.BM) (id int, err error) {
	//log.Println(bs_bm.BsID)

	command := `INSERT INTO bs_bm(
									bs_id, 
									bm_status_id, 
									name, 
									cpf,
                  					cnpj, 
									obs,
									"desc",
									commission_service,
									commission_product,
									commission_supply,
									service_step_time
				) 
				VALUES(:bs_id, 
					   :bm_status_id, 
					   :name, 
					   :cpf,
				       :cnpj,
					   :obs,
					   :desc,
					   :commission_service,
					   :commission_product,
					   :commission_supply,
					   :service_step_time)
					RETURNING ID`

	stmt, err := b.Conn.PrepareNamed(command)

	if err != nil {
		return id, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(bm).Scan(&id)
	return id, err
}

func (b BSBMRepository) UpdateBM(bm entity.BM) error {
	command := `UPDATE bs_bm 
				SET name=:name, 
				    obs=:obs,
					"desc"=:desc,
				    commission_service=:commission_service,
				    commission_product =:commission_product,
					commission_supply =:commission_supply,
					service_step_time=:service_step_time
				WHERE id=:id 
				  AND bs_id=:bs_id`

	stmt, err := b.Conn.PrepareNamed(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bm)
	return b.Error(err)
}

func (b BSBMRepository) ChangeStatus(bsID, bmID, bmStatusID int) error {
	command := `UPDATE bs_bm 
					SET bm_status_id = $3 
					WHERE id=$2 
					  AND bs_id=$1`
	stmt, err := b.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, bmID, bmStatusID)
	return err
}

func (b BSBMRepository) AddBMService(bmID int, service entity.Service) error {
	command := `INSERT INTO bs_bm_service(bs_bm_id,
                          				  bs_service_id,
                          				   duration)
                          			VALUES($1, $2, $3)`
	stmt, err := b.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bmID, service.ID, service.Duration)
	return err
}

func (b BSBMRepository) UpdateBMService(bmID int, service entity.Service) error {
	command := `UPDATE  bs_bm_service 
					SET duration=$3
				WHERE bs_bm_id=$1 
				  AND bs_service_id=$2`
	stmt, err := b.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bmID, service.ID, service.Duration)
	return err
}

func (b BSBMRepository) DeleteBMService(bmID, serviceID int) error {
	// TODO verificar processo de deletar por conta do agendamento
	command := `DELETE FROM bs_bm_service 
						WHERE bs_bm_id = $1 
						  AND bs_service_id = $2`
	stmt, err := b.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bmID, serviceID)
	return err
}

func (b BSBMRepository) AddBMWeekDay(bmID int, day entity.BMWeekDay) error {
	command := `INSERT INTO bs_bm_week_day(
										   bs_bm_id,
                           				   week_day_id,
                           				   start_time,
                           				   end_time )
                           			VALUES($1, $2, $3, $4)`
	stmt, err := b.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bmID, day.ID, day.StartTime, day.EndTime)
	return err
}

func (b BSBMRepository) DeleteBMAllWeekDays(bmID int) error {
	command := `DELETE FROM bs_bm_week_day 
						WHERE bs_bm_id = $1`
	stmt, err := b.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bmID)
	return err
}

func (b BSBMRepository) GetBMS(activePage, itemsPerPage, bsID, bmStatusID int) (bms []entity.BM, totalItems, totalPages int, err error) {

	command := `SELECT 
       				count(*)
                FROM bs_bm bs_bm
                WHERE bs_bm.bs_id=$1 
                  AND bm_status_id=$2`

	err = b.Conn.Get(&totalItems, command, bsID, bmStatusID)
	if err != nil {
		return bms, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT bs_bm.id,
                  	   bs_bm.name,
                  	   bs_bm.cpf,
       				   bs_bm.cnpj,
                  	   bs_bm.obs,
       				   bs_bm."desc"
                FROM bs_bm bs_bm
                WHERE bs_bm.bs_id=$1 
                  AND bm_status_id=$2
                ORDER BY bs_bm.name
	 			OFFSET $3 LIMIT $4`
	if itemsPerPage == 0 {
		itemsPerPage = 1
	}

	err = b.Conn.Select(&bms, command, bsID, bmStatusID, (activePage-1)*itemsPerPage, itemsPerPage)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bms, totalItems, totalPages, err
}

func (b BSBMRepository) GetBMSByName(activePage, itemsPerPage, bsID, bmStatusID int, name string) (bms []entity.BM, totalItems, totalPages int, err error) {

	command := `SELECT 
       				count(*)
                FROM bs_bm bs_bm
                WHERE bs_bm.bs_id=$1 
                  AND bm_status_id=$2
                  AND lower(bs_bm.name) LIKE lower($3)`

	//todo ver questão de busca com acentos.
	err = b.Conn.Get(&totalItems, command, bsID, bmStatusID, "%"+name+"%")
	if err != nil {
		return bms, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT bs_bm.id,
                  	   bs_bm.name,
                  	   bs_bm.cpf,
       				   bs_bm.cnpj,
                  	   bs_bm.obs,
       				   bs_bm."desc"
                FROM bs_bm bs_bm
                WHERE bs_bm.bs_id=$1 
                  AND bm_status_id=$2
                  AND lower(bs_bm.name) LIKE lower($3)
                ORDER BY bs_bm.name
	 			OFFSET $4 LIMIT $5`
	if itemsPerPage == 0 {
		itemsPerPage = 1
	}

	err = b.Conn.Select(&bms, command, bsID, bmStatusID, "%"+name+"%", (activePage-1)*itemsPerPage, itemsPerPage)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bms, totalItems, totalPages, err
}

func (b BSBMRepository) GetBMSByCPF(activePage, itemsPerPage, bsID, bmStatusID int, cpf types.CPF) (bms []entity.BM, totalItems, totalPages int, err error) {

	command := `SELECT 
       				count(*)
                FROM bs_bm bs_bm
                WHERE bs_bm.bs_id=$1 
                  AND bm_status_id=$2
                  AND bs_bm.cpf = $3`

	err = b.Conn.Get(&totalItems, command, bsID, bmStatusID, cpf)
	if err != nil {
		return bms, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT bs_bm.id,
                  	   bs_bm.name,
                  	   bs_bm.cpf,
       				   bs_bm.cnpj,
                  	   bs_bm.obs,
       				   bs_bm."desc"
                FROM bs_bm bs_bm
                WHERE bs_bm.bs_id=$1 
                  AND bm_status_id=$2
                  AND bs_bm.cpf = $3
                ORDER BY bs_bm.name
	 			OFFSET $4 LIMIT $5`
	if itemsPerPage == 0 {
		itemsPerPage = 1
	}

	err = b.Conn.Select(&bms, command, bsID, bmStatusID, cpf, (activePage-1)*itemsPerPage, itemsPerPage)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bms, totalItems, totalPages, err
}

func (b BSBMRepository) GetBmServicesIDs(bmID int) (services []int, err error) {
	command := `SELECT bbs.bs_service_id 
					FROM bs_bm_service bbs 
					WHERE bbs.bs_bm_id=$1`
	err = b.Conn.Select(&services, command, bmID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return services, err
}

func (b BSBMRepository) GetBmServices(bsID, bmID int) (services []entity.Service, err error) {
	command := `SELECT bss.id,
					   bss.name,
					   bss.value,
					   bss.desc,
       				   bbs.duration
					FROM bs_service bss
					INNER JOIN bs_bm_service bbs ON bss.id = bbs.bs_service_id
					WHERE bss.bs_id = $1 AND bbs.bs_bm_id = $2`
	err = b.Conn.Select(&services, command, bsID, bmID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return services, err
}

func (b BSBMRepository) GetBMServices(bmID int) (services []entity.Service, err error) {
	command := `SELECT bse.id,
					   bse.name,
					   bse.value,
					   bse.desc,
       				   bbs.duration
                FROM bs_service bse
                INNER JOIN bs_bm_service bbs ON bse.id = bbs.bs_service_id
                WHERE bbs.bs_bm_id=$1`
	err = b.Conn.Select(&services, command, bmID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return services, err
}

func (b BSBMRepository) GetBMWeekDays(bmID int) (weekDays []entity.BMWeekDay, err error) {
	command := `SELECT wed.id,
					   wed.nome as name,
					   bwd.start_time,
					   bwd.end_time
				FROM weekday wed
				INNER JOIN bs_bm_week_day bwd ON wed.id = bwd.week_day_id
				WHERE bwd.bs_bm_id=$1 
				ORDER BY wed.sequence`
	err = b.Conn.Select(&weekDays, command, bmID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return weekDays, err
}

func (b BSBMRepository) GetBMsByWeekDay(bsID, weekDay int) (weekDays []entity.BMBasic, err error) {
	command := `SELECT bs_bm.id,
					   bs_bm.name
				FROM bs_bm bs_bm 
				INNER JOIN bs_bm_week_day bwd ON bs_bm.id = bwd.bs_bm_id
				WHERE bs_bm.bs_id=$1 
				  AND bwd.week_day_id=$2 
				  AND bs_bm.bm_status_id=1 
				ORDER BY bs_bm.name`
	err = b.Conn.Select(&weekDays, command, bsID, weekDay)
	if err == sql.ErrNoRows {
		err = nil
	}
	return weekDays, err
}
