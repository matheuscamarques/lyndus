package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
)

type ClientCompanyRepository struct {
	conn *sqlx.DB
}

func NewClientCompanyRepository() ClientCompanyRepository {
	return ClientCompanyRepository{
		conn: postgres.DB,
	}
}



func (c ClientCompanyRepository) GetCompanyByID(clientID int) (company entity.Company,err error ){
	command := `SELECT cli.id,
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
					   com.zipcode
					FROM company com
					INNER JOIN client cli on com.id = cli.company_id
					WHERE cli.id = $1 `
	err = c.conn.Get(&company, command, clientID)
	return
}


func (c ClientCompanyRepository) UpdateCompany(company entity.Company) error {
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
                              zipcode=:zipcode
							FROM client
 							WHERE client.id = :id 
 							AND company.id = client.company_id`
	stmt, err := c.conn.PrepareNamed(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(company)
	return err
}




//############
//
//func (c *ClientCompanyRepository) GetCompanyByID(companyID int) error {
//	command := `SELECT com.id,
//					   com.cnpj,
//					   com.company_name,
//					   com.phone,
//					   com.fantasy_name,
//					   com.email,
//					   com.state,
//					   com.city,
//					   com.district,
//					   com.street,
//					   com.number,
//					   com.address_complement,
//					   com.zipcode
//					FROM company com
//					WHERE com.id = $1 `
//	err := postgres.DB.Get(c, command, companyID)
//	return err
//}
//
//
//func (c *Company) GetCompanyByBsID(bsID int) error {
//	command := `SELECT bs.id,
//					   com.cnpj,
//					   com.company_name,
//					   com.phone,
//					   com.fantasy_name,
//					   com.email,
//					   com.state,
//					   com.city,
//					   com.district,
//					   com.street,
//					   com.number,
//					   com.address_complement,
//					   com.zipcode
//					FROM company com
//					INNER JOIN bs  on com.id = bs.company_id
//					WHERE bs.id = $1 `
//	err := postgres.DB.Get(c, command, bsID)
//	return err
//}
//
//
//func (c *Company) UpdateCompany(c *Company) error {
//	command := `UPDATE company SET company_name=:company_name,
//                              phone=:phone,
//                              fantasy_name=:fantasy_name,
//                              email=:email,
//                              state=:state,
//                              city=:city,
//                              district=:district,
//                              street=:street,
//                              number=:number,
//                              address_complement=:address_complement,
//                              zipcode=:zipcode
//                            WHERE id = :id`
//	stmt, err := postgres.DB.PrepareNamed(command)
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//	_, err = stmt.Exec(c)
//	return err
//}
//
//
//
//func (c *Company) UpdateCompanyByBS() error {
//	command := `UPDATE company SET company_name=:company_name,
//                              phone=:phone,
//                              fantasy_name=:fantasy_name,
//                              email=:email,
//                              state=:state,
//                              city=:city,
//                              district=:district,
//                              street=:street,
//                              number=:number,
//                              address_complement=:address_complement,
//                              zipcode=:zipcode
//							FROM company com
//							INNER JOIN bs ON com.id = bs.company_id
//                            WHERE bs.id = :id`
//	stmt, err := postgres.DB.PrepareNamed(command)
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//	_, err = stmt.Exec(c)
//	return err
//}
