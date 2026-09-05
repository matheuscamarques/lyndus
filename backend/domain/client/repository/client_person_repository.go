package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type ClientPersonRepository struct {
	conn *sqlx.DB
}

func NewClientPersonRepository() ClientPersonRepository {
	return ClientPersonRepository{
		conn: postgres.DB,
	}
}

func (p ClientPersonRepository) CreatePersonBasic(person *entity.PersonBasic) error {
	command := `INSERT INTO person(name, cpf) VALUES($1, $2) RETURNING ID`
	stmt, err := p.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(person.Name, person.CPF).Scan(&person.ID)
	return err
}

func (p ClientPersonRepository) GetPersonIDByCPF(person *entity.PersonBasic) error {
	command := `SELECT per.id FROM person per WHERE per.cpf=$1`
	err := p.conn.Get(person, command, person.CPF)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
