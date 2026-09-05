package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type BSCompanyRepository struct {
	conn *sqlx.DB
}

func NewBSCompanyRepository() BSCompanyRepository {
	return BSCompanyRepository{
		conn: postgres.DB,
	}
}

func (b BSCompanyRepository) GetLocation(bsID int) (location string, err error) {
	command := `SELECT location
					FROM bs 
					WHERE bs.id = $1 `
	err = b.conn.Get(&location, command, bsID)
	return location, err
}

func (b BSCompanyRepository) GetCompanyByBsID(bsID int) (company entity.Company, err error) {
	command := `SELECT bs.id,
					   com.cnpj,
					   com.company_name,
					   com.phone,
					   com.fantasy_name,
					   com.email,
					   com.state,
					   com.city,
					   com.district,
					   com.street,
					   com.number,
					   com.address_complement,
					   com.zipcode,
					   com.desc
					FROM company com
					INNER JOIN bs  on com.id = bs.company_id
					WHERE bs.id = $1 `
	err = b.conn.Get(&company, command, bsID)
	return company, err
}

func (b BSCompanyRepository) UpdateCompanyByBS(company entity.Company) error {
	command := `UPDATE company SET company_name=:company_name,
                              phone=:phone,
                              fantasy_name=:fantasy_name,
                              email=:email,
                              state=:state,
                              city=:city,
                              district=:district,
                              street=:street,
                              number=:number,
                              address_complement=:address_complement,
                              zipcode=:zipcode,
							  "desc"=:desc
							FROM bs
 							WHERE bs.id = :id 
 							AND company.id = bs.company_id`
	stmt, err := b.conn.PrepareNamed(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(company)
	return err
}

func (b BSCompanyRepository) GetBSWeekDays(bsID int) (weekDays []entity.BMWeekDay, err error) {
	command := `SELECT wed.id,
					   wed.nome as name,
					   bwd.start_time,
					   bwd.end_time
				FROM weekday wed
				INNER JOIN bs_week_day bwd ON wed.id = bwd.week_day_id
				WHERE bwd.bs_id=$1 ORDER BY wed.sequence`
	err = b.conn.Select(&weekDays, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return weekDays, err
}

func (b BSCompanyRepository) DeleteBSAllWeekDays(bsID int) error {
	command := `DELETE FROM bs_week_day WHERE bs_id = $1`
	stmt, err := b.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID)
	return err
}

func (b BSCompanyRepository) AddBSWeekDay(bsID int, day entity.BMWeekDay) error {
	command := `INSERT INTO bs_week_day(bs_id, week_day_id, start_time, end_time ) VALUES($1, $2, $3, $4)`
	stmt, err := b.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, day.ID, day.StartTime, day.EndTime)
	return err
}
