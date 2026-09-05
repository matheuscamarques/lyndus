package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type BSServiceRepository struct {
	contracts.BSServiceRepositoryInterface
	conn *sqlx.DB
}

func NewBSServiceRepository() *BSServiceRepository {
	return &BSServiceRepository{
		conn: postgres.DB,
	}
}

func (s BSServiceRepository) GetBSServiceCategories() (serviceCategories []entity.ServiceCategory, err error) {
	command := `SELECT id,
       				   "desc"
       				FROM bs_service_category
       				order by id`
	err = s.conn.Select(&serviceCategories, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return serviceCategories, err
}

func (s BSServiceRepository) GetServiceID(bsID, serviceID int) (n int, err error) {
	command := `SELECT bss.id
					FROM bs_service bss
					WHERE bss.bs_id=$1 AND bss.id=$2`
	err = s.conn.Get(&n, command, bsID, serviceID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return n, err
}

func (s BSServiceRepository) GetService(bsID, serviceID int) (service entity.Service, err error) {
	command := `SELECT bse.id,
                  	   bse.name,
                  	   bse.value,
                	   bse."desc",
       				   bse.bs_service_category_id,
       				   bsc."desc" as category_name
                FROM bs_service bse
				INNER JOIN bs_service_category bsc ON bse.bs_service_category_id = bsc.id
                WHERE bs_id=$1 AND bse.id=$2`
	err = s.conn.Get(&service, command, bsID, serviceID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return service, err
}

func (s BSServiceRepository) GetBSServices(bsID int) (services []entity.Service, err error) {
	command := `SELECT bse.id,
       				   bse.name,
       				   bse.value,
       				   bse."desc",
					   bsc."desc" as category_name
       				FROM bs_service bse
					INNER JOIN bs_service_category bsc on bsc.id = bse.bs_service_category_id
					WHERE bs_id=$1 `
	err = s.conn.Select(&services, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return services, err
}

func (s BSServiceRepository) UpdateService(service entity.Service) error {
	//log.Println(bs_service)
	command := `UPDATE bs_service SET name=$3,
                      				  "desc"=$4,
                      				  value=$5,
									  bs_service_category_id=$6
								WHERE id=$1 AND bs_id=$2`
	stmt, err := s.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(service.ID, service.BsID, service.Name, service.Desc, service.Value, service.CategoryID)
	return err
}

func (s BSServiceRepository) GetBSServicesByName(bsID int, name string) (services []entity.Service, err error) {
	command := `SELECT bse.id,
                  	   bse.name,
                  	   bse.value,
                  	   bse."desc",
       			  	   bsc."desc" as category_name
                FROM bs_service bse
				INNER JOIN bs_service_category bsc on bsc.id = bse.bs_service_category_id
                WHERE bs_id=$1 AND lower(name) LIKE lower($2)`
	err = s.conn.Select(&services, command, bsID, name)
	if err == sql.ErrNoRows {
		err = nil
	}
	return services, err
}

func (s BSServiceRepository) CreateService(service entity.Service) (id int, err error) {
	command := `INSERT INTO bs_service(bs_id, 
                       				   name,
                       				   value,
                       				   "desc",
                       				   bs_service_category_id)
                       				VALUES($1, $2, $3, $4, $5)
                       				RETURNING ID`
	stmt, err := s.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(service.BsID,
		service.Name,
		service.Value,
		service.Desc,
		service.CategoryID).Scan(&id)
	return id, err
}

func (s BSServiceRepository) AddBMService(bmID int, service entity.Service) error {
	command := `INSERT INTO bs_bm_service(bs_bm_id,
                          				  bs_service_id,
                          				  duration
                          				  )
                          			VALUES($1, $2 , $3)`
	stmt, err := s.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bmID, service.ID, service.Duration)
	return err
}

func (s BSServiceRepository) DeleteBMService(bmID, serviceID int) error {
	// TODO verificar processo de deletar por conta do agendamento
	command := `DELETE FROM bs_bm_service WHERE bs_bm_id = $1 AND bs_service_id = $2`
	stmt, err := s.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bmID, serviceID)
	return err
}

func (s BSServiceRepository) GetBankDetails(bsID int) (details entity.BSBankDetails, err error) {
	command := `SELECT  id, 
       					bs_id, 
       					bank, 
       					branch, 
       			bank_account FROM bs_bank_details WHERE bs_id = $1`
	err = s.conn.Get(&details, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return details, err
}

func (s BSServiceRepository) UpdateBankDetails(bank entity.BSBankDetails) error {
	command := `UPDATE bs_bank_details SET 
                           				bank = $2,
                           				agency = $3,
										bank_account = $4
								WHERE  bs_id=$1`
	stmt, err := s.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		bank.BsID,
		bank.Bank,
		bank.Branch,
		bank.BankAccount,
	)
	return err
}

func (s BSServiceRepository) CreateBankDetails(bank entity.BSBankDetails) error {
	command := `INSERT INTO bs_bank_details(bs_id, bank, agency, bank_account)
                       				VALUES($1, $2, $3, $4)
                       				RETURNING ID`
	stmt, err := s.conn.Prepare(command)
	var id int
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		bank.BsID,
		bank.Bank,
		bank.Branch,
		bank.BankAccount,
	).Scan(&id)
	return err
}
