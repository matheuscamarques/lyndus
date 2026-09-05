package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/jmoiron/sqlx"
)

type BonderCompanyRepository struct {
	conn *sqlx.DB
}

func NewBonderCompanyRepository() BonderCompanyRepository {
	return BonderCompanyRepository{
		conn: postgres.DB,
	}
}
func (bcr BonderCompanyRepository) GetIdByCNPJ(cnpj types.CNPJ) (companyID int, err error) {
	command := "SELECT id FROM company WHERE cnpj=$1 LIMIT 1"
	err = bcr.conn.Get(&companyID, command, cnpj)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return 0, err
	}
	return companyID, err
}

func (bcr BonderCompanyRepository) Create(company *entity.BS) (companyID int, err error) {
	command := `INSERT INTO company(
								cnpj, 
								company_name,
								phone,
								fantasy_name,
								email,
								state,
								city,
								district,
								street,
								number,
								address_complement,
								zipcode
							)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id
		 `
	stmt, err := bcr.conn.Prepare(command)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(
		company.CNPJ,
		company.CompanyName,
		company.Phone,
		company.FantasyName,
		company.Email,
		company.State,
		company.City,
		company.District,
		company.Street,
		company.Number,
		company.AddressComplement,
		company.Zipcode,
	).Scan(&id)
	return id, err
}

func (bcr BonderCompanyRepository) Update(company *entity.BS) error {
	command := `UPDATE company SET 
              				company_name=$2,
							phone=$3,
							fantasy_name=$4,
							email=$5,
							state=$6,
							city=$7,
							district=$8,
							street=$9,
							number=$10,
							address_complement=$11,
							zipcode=$12 
			   WHERE id = $1`
	_, err := bcr.conn.Exec(command, company.ID,
		company.CompanyName,
		company.Phone,
		company.FantasyName,
		company.Email,
		company.State,
		company.City,
		company.District,
		company.Street,
		company.Number,
		company.AddressComplement,
		company.Zipcode)
	return err
}
