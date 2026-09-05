package repository

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type AppUserCreateRepository struct {
	contracts.AppUserCreateRepositoryInterface
	conn *sqlx.DB
}

func NewAppUserCreateRepository() *AppUserCreateRepository {
	return &AppUserCreateRepository{
		conn: postgres.DB,
	}
}

//CreatePerson func to create banner
func (aucr AppUserCreateRepository) CreatePerson(auc entity.AppUserCreate) (id int, err error) {
	command := `INSERT INTO bs_person(gender_id, 
                   				   name, 
                   				   cpf, 
                   				   birthdate, 
                   				   phone, 
                   				   cellphone, 
                   				   email,
            					   state, 
                   				   city, 
                   				   district, 
                   				   street, 
                   				   number, 
                   				   address_complement, 
                   				   zipcode)
            VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING ID`
	stmt, err := aucr.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(auc.GenderID, auc.Name, auc.CPF, auc.Birthdate, auc.Phone, auc.Cellphone, auc.Email, auc.State,
		auc.City, auc.District, auc.Street, auc.Number, auc.AddressComplement, auc.Zipcode).Scan(&id)
	return id, err
}

func (aucr AppUserCreateRepository) UpdatePerson(auc entity.AppUserCreate) error {
	command := `UPDATE bs_person SET gender_id=$2,
								  name=$3,
								  birthdate=$4,
								  phone=$5,
								  cellphone=$6,
								  email=$7,
								  state=$8,
								  city=$9,
								  district=$10,
								  street=$11,
								  number=$12,
								  address_complement=$13,
								  zipcode=$14  
                        WHERE id=$1`
	stmt, err := aucr.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(auc.ID, auc.GenderID, auc.Name, auc.Birthdate, auc.Phone, auc.Cellphone, auc.Email, auc.State,
		auc.City, auc.District, auc.Street, auc.Number, auc.AddressComplement, auc.Zipcode)
	return err
}

func (aucr AppUserCreateRepository) GetPersonIDByCPF(a *entity.AppUserCreate) (n int, err error) {
	command := `SELECT per.id FROM bs_person per WHERE per.cpf=$1`
	err = aucr.conn.Get(&n, command, a.CPF)
	if err == sql.ErrNoRows {
		err = nil
	}
	return n, err
}

func (aucr AppUserCreateRepository) GetAppUserIDByCPF(a *entity.AppUserCreate) (n int, err error) {
	command := `SELECT app.id  FROM app_user app WHERE app.cpf = $1`
	err = aucr.conn.Get(&n, command, a.CPF)
	if err == sql.ErrNoRows {
		err = nil
	}
	return n, err
}
