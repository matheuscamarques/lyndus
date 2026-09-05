package repository

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"
	"bitbucket.org/lyndus/backend/infra/types"
	"bitbucket.org/lyndus/backend/internal/constants"
	"database/sql"
	"github.com/shopspring/decimal"
	"log"
)

type AppUserRepository struct {
	db.Connector
}

func NewAppUserRepository() (repo AppUserRepository) {
	repo.Table = "app_user"
	repo.UsePostgres()
	return repo
}

func (aur AppUserRepository) GetAppUserByID(appUserID int) (appUsers entity.AppUser, err error) {
	command := `SELECT aus.id,
       			       per.name, 
					   per.cpf, 
					   per.birthdate, 
					   per.phone, 
					   per.cellphone,
       				   per.email,
					   per.state, 
					   per.city, 
					   per.district, 
					   per.street, 
					   per.number, 
					   per.address_complement, 
					   per.zipcode, 
       				   aub.value as faz_bem,
					   aulb.value as lyndus_box
       			FROM app_user as aus 
				INNER JOIN person as per ON aus.person_id = per.id
				LEFT JOIN app_user_benefit aub on aus.id = aub.app_user_id
				LEFT JOIN app_user_lyndus_box aulb on aus.id = aulb.app_user_id
					WHERE aus.id = $1 
					ORDER BY per.name `

	err = aur.Conn.Get(&appUsers, command, appUserID)
	return appUsers, err
}

func (aur AppUserRepository) GetAppUserByCPF(cpf types.CPF, criteria criteria.Criteria) (appusers []aggregate.AppUserAggreate, err error) {
	command := `SELECT aus.id,
       			       per.name, 
					   per.cpf, 
					   per.birthdate, 
					   per.phone, 
					   per.cellphone,
       				   per.email
       			FROM app_user as aus 
				INNER JOIN person as per ON aus.person_id = per.id
					WHERE aus.cpf = $1 
					ORDER BY per.name `

	log.Println("cmd1", command)
	command = criteria.Query(command)
	log.Println("cmd2", command)
	err = aur.Conn.Select(&appusers, command, cpf)
	return appusers, err
}

func (aur AppUserRepository) GetAppUsers(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (appUsers []entity.AppUserPerson, totalItems, totalPages int, err error) {

	command := `SELECT count(*)
					FROM app_user apu
					INNER JOIN person per ON apu.person_id = per.id `

	if search != "" {
		command = command + ` WHERE unaccent(per.name) ILIKE $1  `
		err = aur.Conn.Get(&totalItems, command, "%"+search+"%")
	} else {
		err = aur.Conn.Get(&totalItems, command)
	}

	if err != nil {
		return appUsers, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT apu.id,
       			       per.name, 
					   per.cpf, 
					   per.birthdate, 
					   per.phone, 
					   per.cellphone,
       				   per.email
       			FROM app_user apu
					INNER JOIN person per ON apu.person_id = per.id
					`

	if search != "" {
		command = command + ` WHERE unaccent(per.name) ILIKE $3 `
	}

	if orderBy == "name" {
		command = command + ` ORDER BY per.name `
	} else if orderBy == "cpf" {
		command = command + ` ORDER BY per.cpf `
	} else if orderBy == "birthdate" {
		command = command + ` ORDER BY per.birthdate `
	} else if orderBy == "phone" {
		command = command + ` ORDER BY per.phone `
	} else if orderBy == "email" {
		command = command + ` ORDER BY per.email `
	} else {
		command = command + ` ORDER BY per.name `
	}

	if sortDesc {
		command = command + " DESC "
	}
	command = command + ` OFFSET $1 LIMIT $2`

	if search != "" {
		err = aur.Conn.Select(&appUsers, command, (activePage-1)*itemsPerPage, itemsPerPage, "%"+search+"%")
	} else {
		err = aur.Conn.Select(&appUsers, command, (activePage-1)*itemsPerPage, itemsPerPage)
	}
	if err == sql.ErrNoRows {
		err = nil
	}

	return appUsers, totalItems, totalPages, err
}

func (aur AppUserRepository) GetAppUserByName(name string, criteria criteria.Criteria) (appusers []aggregate.AppUserAggreate, err error) {
	command := `SELECT aus.id,
       			       per.name, 
					   per.cpf, 
					   per.birthdate, 
					   per.phone, 
					   per.cellphone,
       				   per.email
       			FROM app_user as aus 
				INNER JOIN person as per ON aus.person_id = per.id
					WHERE lower(per.name) LIKE lower($1)
					ORDER BY per.name `

	command = criteria.Query(command)
	err = aur.Conn.Select(&appusers, command, "%"+name+"%")
	return appusers, err
}

func (aur AppUserRepository) GetBalance(appUserID int) (bonus entity.Bonus, err error) {
	command := `SELECT aub.id,
       				   aub.value as balance
       			FROM  app_user_benefit aub 
					WHERE aub.app_user_id = $1`

	err = aur.Conn.Get(&bonus, command, appUserID)
	return bonus, err
}

//func (aur AppUserRepository) InsertStatement(statement entity.AppUserStatement) (int, error) {
//	command := `INSERT INTO app_user_statement(app_user_balance_id,
//                               				   app_user_id,
//                               				   value,
//                               				   statements_id,
//                               				   "desc",
//                               				   date_time)
//						VALUES($1,$2,$3,$4,$5,now()) RETURNING id`
//	stmt, err := aur.Conn.Prepare(command)
//	if err != nil {
//		return 0, err
//	}
//	defer stmt.Close()
//	var id int
//	err = stmt.QueryRow(statement.AppUserBalanceID,
//		statement.AppUserID,
//		statement.Value,
//		statement.StatementID,
//		statement.Desc).Scan(&id)
//	return id, err
//}

func (aur AppUserRepository) InsertBalance(bonus entity.Bonus) (int, error) {
	command := "INSERT INTO app_user_benefit(app_user_id, value, last_update) VALUES($1,$2,now()) RETURNING id"
	stmt, err := aur.Conn.Prepare(command)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(bonus.ID, bonus.Balance).Scan(&id)
	return id, err
}

func (aur AppUserRepository) InsertBalanceHistory(appUserID int, oldValue, newValue decimal.Decimal) (int, error) {
	command := "INSERT INTO public.app_user_benefit_history(app_user_id, value_old, value_new) VALUES($1,$2,$3) RETURNING id"
	stmt, err := aur.Conn.Prepare(command)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(appUserID, oldValue, newValue).Scan(&id)
	return id, err
}

func (aur AppUserRepository) UpdateBalance(appUserID int, value decimal.Decimal) error {
	command := `UPDATE app_user_benefit SET value=$2, last_update=now() WHERE app_user_id=$1`
	stmt, err := aur.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(appUserID, value)
	return err
}

func (aur AppUserRepository) UpdateLyndusBox(bonus entity.Bonus) error {
	paymentCode := "bonus-bonder"
	sqlCommand := `SELECT "value" FROM app_user_lyndus_box WHERE app_user_id=$1`
	var oldValue decimal.Decimal
	var insert bool

	err := aur.Conn.Get(&oldValue, sqlCommand, bonus.ID)

	if err == sql.ErrNoRows {
		insert = true
		err = nil
	}
	if err != nil {
		return err
	}

	tx, err := aur.Conn.Begin()
	if err != nil {
		return err
	}

	value := bonus.Balance.Add(oldValue)
	if insert {
		sqlCommand = `INSERT INTO app_user_lyndus_box (app_user_id,
													   "value",
													   last_update) 
												VALUES($1, $2, now())`
	} else {
		sqlCommand = `UPDATE app_user_lyndus_box SET "value"=$2, last_update=now() WHERE app_user_id=$1`
	}

	_, err = tx.Exec(sqlCommand, bonus.ID, value)
	if err != nil {
		tx.Rollback()
		return err
	}

	//sqlCommand = `INSERT INTO app_user_statement(app_user_id,
	//										     value,
	//										     "desc",
	//										     statements_id,
	//										     statements_methods_id)
	//							VALUES($1, $2, $3, $4, $5) RETURNING ID `
	//
	var userStatementsID int
	//err = tx.QueryRow(sqlCommand,
	//	appUserID,
	//	boxValue,
	//	constants.BYLyndusBoxText,
	//	constants.LyndusPay,
	//	constants.DEBIT).Scan(&userStatementsID)
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}

	sqlCommand = `INSERT INTO app_user_statement(app_user_id,
											     value,
											     "desc",
											     statements_id,
											     statements_methods_id)
								VALUES($1, $2, $3, $4, $5) RETURNING ID`

	err = tx.QueryRow(sqlCommand,
		bonus.ID,
		bonus.Balance,
		constants.AddBonusLyndusBoxText,
		constants.LyndusBox,
		constants.CREDIT).Scan(&userStatementsID)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlCommand = `INSERT INTO app_user_lyndus_box_history(app_user_id,
														  value_old,
														  value_new,
														  payment_code,
														  app_user_statement_id)
												VALUES($1,$2,$3,$4,$5)`

	_, err = tx.Exec(sqlCommand, bonus.ID, oldValue, value, paymentCode, userStatementsID)
	if err != nil {
		tx.Commit()
		return err
	}

	err = tx.Commit()
	return err
}
