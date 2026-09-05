package repository

import (
	"database/sql"

	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/jmoiron/sqlx"
)

type BSPersonRepository struct {
	contracts.BSPersonRepositoryInterface
	conn *sqlx.DB
}

func NewBSPersonRepository() BSPersonRepository {
	return BSPersonRepository{
		conn: postgres.DB,
	}
}

func (p BSPersonRepository) FindPerson(bsID int, cpf *types.CPF, phone string) (person entity.Person, err error) {
	command := `SELECT id,
					   app_user_id,
					   name,
					   cpf,
					   phone,
					   birthdate
					FROM bs_person
					WHERE bs_id=$1 AND (cpf=$2 OR phone=$3)`
	err = p.conn.Get(&person, command, bsID, cpf, phone)
	if err == sql.ErrNoRows {
		err = nil
	}
	return person, err
}

func (p BSPersonRepository) GetPersonID(bsID, personID int) (person entity.Person, err error) {
	command := `SELECT id,
					   app_user_id
					FROM bs_person
					WHERE bs_id=$1 AND id=$2`
	err = p.conn.Get(&person, command, bsID, personID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return person, err
}

func (p BSPersonRepository) GetPerson(bsID, personID int) (person entity.Person, err error) {
	command := `SELECT id,
					   app_user_id,
					   name,
					   cpf,
					   phone,
					   birthdate 
					FROM bs_person
					WHERE bs_id=$1 AND id=$2`
	err = p.conn.Get(&person, command, bsID, personID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return person, err
}

func (p BSPersonRepository) GetPersons(bsID int) (persons []entity.Person, err error) {
	command := `SELECT per.id,
					   per.name,
					   per.cpf,
					   per.phone,
					   per.birthdate
                FROM bs_person per
                WHERE per.bs_id=$1 ORDER BY per.name`
	err = p.conn.Select(&persons, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return persons, err
}

func (p BSPersonRepository) GetPersonsByCPF(bsID int, cpf types.CPF) (persons []entity.Person, err error) {
	command := `SELECT per.id,
                  	   per.name,
					   per.cpf,
                  	   per.phone,
       				   per.birthdate
				FROM bs_person per
                WHERE per.bs_id = $1 AND per.cpf = $2`
	err = p.conn.Select(&persons, command, bsID, cpf)
	if err == sql.ErrNoRows {
		err = nil
	}
	return persons, err
}

func (p BSPersonRepository) GetPersonsByPhone(bsID int, phone string) (persons []entity.Person, err error) {
	command := `SELECT per.id,
                  	   per.name,
                  	   per.cpf,
                  	   per.phone,
                  	   per.birthdate
                FROM bs_person per
                WHERE per.bs_id=$1 AND lower(per.phone) LIKE lower($2)`
	err = p.conn.Select(&persons, command, bsID, phone)
	if err == sql.ErrNoRows {
		err = nil
	}
	return persons, err
}

func (p BSPersonRepository) GetPersonsByName(bsID int, name string) (persons []entity.Person, err error) {
	command := `SELECT per.id,
                  	   per.name,
                  	   per.cpf,
                  	   per.phone,
                  	   per.birthdate
                FROM bs_person per
                WHERE per.bs_id=$1 AND lower(per.name) LIKE lower($2) ORDER BY per.name`
	err = p.conn.Select(&persons, command, bsID, name+"%")
	if err == sql.ErrNoRows {
		err = nil
	}
	return persons, err
}

func (p BSPersonRepository) CreatePerson(person entity.Person) (id int, err error) {
	command := `INSERT INTO bs_person(bs_id,
                      				  app_user_id,
                      				  name,
                      				  cpf,
                      				  phone,
                      				  birthdate) 
            		VALUES($1, $2, $3, $4, $5, $6) RETURNING ID`
	stmt, err := p.conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(person.BsID, person.AppUserID, person.Name, person.CPF, person.Phone, person.Birthdate).Scan(&id)
	return id, err
}

func (p BSPersonRepository) UpdatePerson(person entity.Person) error {
	command := `UPDATE bs_person SET app_user_id=$3,
                     				 name=$4,
                     				 cpf=$5,
                     				 phone=$6,
                     				 birthdate=$7 
						WHERE bs_id=$1 AND id=$2`
	stmt, err := p.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(person.BsID, person.ID, person.AppUserID, person.Name, person.CPF, person.Phone, person.Birthdate)
	return err
}

func (p BSPersonRepository) UpdatePersonCPF(bsID, personID int, cpf types.CPF) error {
	command := `UPDATE bs_person SET cpf=$3 
								WHERE bs_id=$1 AND id=$2`
	stmt, err := p.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, personID, cpf)
	return err
}

func (p BSPersonRepository) UpdatePersonCPFAndAppUser(bsID, personID int, cpf types.CPF, appUserID int) error {
	command := `UPDATE bs_person SET cpf=$3, 
                     			     app_user_id=$4 
								 WHERE bs_id=$1
									AND id=$2`
	stmt, err := p.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bsID, personID, cpf, appUserID)
	return err
}
