package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/infra/db"
	"database/sql"
)

type ClientBenefitRepository struct {
	db.Connector
}

func NewClientBenefitRepository() (repo ClientBenefitRepository) {
	repo.UsePostgres()
	repo.Table = "client_benefit"
	return repo
}
func (br ClientBenefitRepository) CreateBenefitUserCategory(buc *entity.BenefitUserCategory) (id int, err error) {
	command := `INSERT INTO client_benefit_user_category(client_benefit_user_id, client_category_id, name, value)
					VALUES($1, $2, $3, $4) RETURNING ID`
	stmt, err := br.Conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(buc.BenefitUserID, buc.CategoryID, buc.Name, buc.Value).Scan(&id)
	return id, err
}

func (br ClientBenefitRepository) CreateBenefit(b *entity.Benefit) error {
	command := `INSERT INTO client_benefit(client_id, benefit_status_id, value, "desc")
				 VALUES($1, $2, $3, $4) RETURNING ID`
	stmt, err := br.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(b.ClientID, b.BenefitStatusID, b.Value, b.Desc).Scan(&b.ID)
	return err
}

func (br ClientBenefitRepository) CreateBenefitUser(bu entity.BenefitUser) (id int, err error) {
	command := `INSERT INTO client_benefit_user(client_id, client_benefit_id, app_user_id, value, client_employee_id)
					VALUES($1, $2, $3, $4, $5) RETURNING ID`
	stmt, err := br.Conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(bu.ClientID, bu.BenefitID, bu.AppUserID, bu.Value, bu.EmployeeID).Scan(&id)
	return id, err
}

func (br ClientBenefitRepository) UpdateBenefitValue(b *entity.Benefit) error {
	command := `UPDATE client_benefit SET value=$3 WHERE client_id=$1 AND id=$2`
	stmt, err := br.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(b.ClientID, b.ID, b.Value)
	return err
}

func (br ClientBenefitRepository) UpdateBenefitUserValue(bu *entity.BenefitUser) error {
	command := `UPDATE client_benefit_user SET value=$3 WHERE client_id=$1 AND id=$2`
	stmt, err := br.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bu.ClientID, bu.ID, bu.Value)
	return err

}

func (br ClientBenefitRepository) UpdateBenefitUserAdditional(bu *entity.BenefitUser) error {
	command := `UPDATE client_benefit_user SET additional_value=$4,
			                                    additional_reason=$5 
				WHERE client_id=$1 AND client_benefit_id=$2 AND id=$3`
	stmt, err := br.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bu.ClientID, bu.BenefitID, bu.ID, bu.AdditionalValue, bu.AdditionalReason)
	return err
}

func (br ClientBenefitRepository) GetBenefitStatus(b *entity.Benefit) (status int, err error) {
	command := `SELECT 
					  cbe.benefit_status_id
					FROM client_benefit cbe
					WHERE cbe.client_id = $1 AND cbe.id= $2 `
	err = br.Conn.Get(&status, command, b.ClientID, b.ID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return status, err
}

func (br ClientBenefitRepository) GetBenefit(b entity.Benefit) (benefitResp entity.BenefitRespFull, err error) {
	command := `SELECT cbe.id,
					  bs.desc_pt as benefit_status,
       				  bs.id as benefit_status_id,
					  cbe.value,
					  cbe.desc,
					  cbe.created_at as "date"
					FROM client_benefit cbe
					INNER JOIN benefit_status bs on bs.id = cbe.benefit_status_id
					WHERE cbe.client_id = $1 AND cbe.id= $2 `
	err = br.Conn.Get(&benefitResp, command, b.ClientID, b.ID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return benefitResp, err
}

func (br ClientBenefitRepository) UpdateBenefitDesc(b entity.Benefit) error {
	command := `UPDATE client_benefit SET "desc"=$3 WHERE client_id=$1 AND id=$2`
	stmt, err := br.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(b.ClientID, b.ID, b.Desc)
	return err
}

func (br ClientBenefitRepository) UpdateBenefitStatus(b entity.Benefit) error {
	command := `UPDATE client_benefit SET benefit_status_id=$3 WHERE client_id=$1 AND id=$2`
	stmt, err := br.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(b.ClientID, b.ID, b.BenefitStatusID)
	return err

}

func (br ClientBenefitRepository) GetBenefitByStatus(b entity.Benefit) error {
	command := `SELECT cbe.id,
					   cbe.benefit_status_id,
					   cbe.value,
					   cbe.desc,
					   cbe.created_at as "date"
					FROM client_benefit cbe
					WHERE cbe.client_id =$1 AND cbe.id = $2 AND cbe.benefit_status_id = $3`
	err := br.Conn.Get(b, command, b.ClientID, b.ID, b.BenefitStatusID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (br ClientBenefitRepository) GetBenefits(clientID int) (benefitsResp []entity.BenefitResp, err error) {
	command := `SELECT cbe.id,
					   bs.desc_pt as benefit_status,
					   cbe.value,
					   cbe.desc as "desc",
					   cbe.created_at as "date"
					FROM client_benefit cbe
					INNER JOIN benefit_status bs on cbe.benefit_status_id = bs.id
					WHERE client_id= $1
					ORDER BY created_at DESC`
	err = br.Conn.Select(&benefitsResp, command, clientID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return benefitsResp, err
}

// func (a *Authentication) get_benefits( companyID int):
//     return await db.fetch_rows("""SELECT cbe.id, bst.desc as benefit_status, cbe.value, cbe.desc as desc, cbe.created_at as "date" FROM company_benefit cbe
//         INNER JOIN benefit_status bst ON cbe.benefit_status_id = bst.id
//         WHERE company_id=$1  """, company_id)

func (br ClientBenefitRepository) GetBenefitUser(b entity.Benefit) (benefitUser entity.BenefitUserResp, err error) {
	command := `SELECT id,
					  value,
					  additional_value,
					  additional_reason
					FROM client_benefit_user
					WHERE client_id = $1 AND id = $2 `
	err = br.Conn.Get(&benefitUser, command, b.ClientID, b.ID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return benefitUser, err
}

func (br ClientBenefitRepository) GetBenefitsUsers(b entity.Benefit) (benefitsUser []entity.BenefitUserResp, err error) {
	command := `SELECT cbu.id,
					   cbu.value,
					   cbu.additional_value,
					   cbu.additional_reason,
					   cem.name,
					   cem.cpf
					FROM client_benefit_user cbu
					INNER JOIN client_employee cem ON cem.id = cbu.client_employee_id
					WHERE cbu.client_id = $1 AND cbu.client_benefit_id= $2
					ORDER BY cem.name `
	err = br.Conn.Select(&benefitsUser, command, b.ClientID, b.ID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return benefitsUser, err
}

func (br ClientBenefitRepository) GetBenefitsUserCategory(bu entity.BenefitUserResp) (categories []string, err error) {
	command := `SELECT 
                  cbuc.name as category 
                FROM client_benefit_user cbu  
                INNER JOIN client_benefit_user_category cbuc ON cbu.id =  cbuc.client_benefit_user_id 
	            WHERE cbu.id=$1`
	err = br.Conn.Select(&categories, command, bu.ID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return categories, err
}

func (br ClientBenefitRepository) GetAll(activePage, itemsPerPage, clientID, statusID int, search, orderBy string, sortDesc bool) (benefits []entity.BenefitResp, totalItems, totalPages int, err error) {
	command := `SELECT count(*)
					FROM client_benefit cbe
					INNER JOIN benefit_status bes on cbe.benefit_status_id = bes.id
					WHERE cbe.client_id= $1`

	if search != "" {
		if statusID != 0 {
			command = command + ` AND cbe.benefit_status_id = $2  AND unaccent(cbe.desc) ILIKE $3 `
			err = br.Conn.Get(&totalItems, command, clientID, statusID, "%"+search+"%")
		} else {
			command = command + ` AND unaccent(cbe.desc) ILIKE $2 `
			err = br.Conn.Get(&totalItems, command, clientID, "%"+search+"%")
		}

	} else {
		if statusID != 0 {
			command = command + ` AND cbe.benefit_status_id = $2 `
			err = br.Conn.Get(&totalItems, command, clientID, statusID)
		} else {
			err = br.Conn.Get(&totalItems, command, clientID)
		}
	}

	if err != nil {
		return benefits, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT cbe.id,
					   bes.desc_pt as benefit_status,
       				   cbe.benefit_status_id,	
					   cbe.value,
					   cbe.desc as "desc",
					   cbe.created_at as "date"
					FROM client_benefit cbe
					INNER JOIN benefit_status bes on cbe.benefit_status_id = bes.id
					WHERE cbe.client_id= $3
					 `

	if search != "" {
		if statusID != 0 {
			command += ` AND cbe.benefit_status_id = $4 AND unaccent(cbe.desc) ILIKE $5`
		} else {
			command += ` AND unaccent(cbe.desc) ILIKE $4`
		}
	} else if statusID != 0 {
		command += ` AND cbe.benefit_status_id = $4 `

	}

	if orderBy == "desc" {
		command = command + ` ORDER BY cbe.desc `
	} else if orderBy == "date" {
		command = command + ` ORDER BY cbe.created_at `
	} else if orderBy == "value" {
		command = command + ` ORDER BY cbe.value `
	} else {
		command = command + ` ORDER BY cbe.desc `
	}

	if sortDesc {
		command = command + " DESC "
	}
	command = command + ` OFFSET $1 LIMIT $2`

	if search != "" {
		if statusID != 0 {
			err = br.Conn.Select(&benefits, command, (activePage-1)*itemsPerPage, itemsPerPage, clientID, statusID, "%"+search+"%")
		} else {
			err = br.Conn.Select(&benefits, command, (activePage-1)*itemsPerPage, itemsPerPage, clientID, "%"+search+"%")
		}
	} else {
		if statusID != 0 {
			err = br.Conn.Select(&benefits, command, (activePage-1)*itemsPerPage, itemsPerPage, clientID, statusID)
		} else {
			err = br.Conn.Select(&benefits, command, (activePage-1)*itemsPerPage, itemsPerPage, clientID)
		}
	}
	if err == sql.ErrNoRows {
		err = nil
	}

	return benefits, totalItems, totalPages, err
}

//// multicategory
// func (a *Authentication) get_benefits_users( companyID int, benefitID int):
//     return await db.fetch_rows("""SELECT  cbu.id, cbu.value, cbu.additional_value, cbu.additional_reason, cem.name, cem.cpf
//         FROM company_benefit_user cbu
//         INNER JOIN company_employee cem ON cem.id = cbu.company_employee_id
//         where cbu.company_id=$1 AND cbu.company_benefit_id=$2 ORDER BY cem.name  """, company_id, benefit_id)

// func (a *Authentication) get_benefits_user_category( benefit_userID int):
//     return await db.fetch_rows("""SELECT name, value  FROM company_benefit_user_category cbuc
//         WHERE cbuc.company_category_user_id """, benefit_user_id)

func (br ClientBenefitRepository) GetBenefitsUsersValues(bu entity.BenefitUser) (benefitsUser []entity.BenefitUser, err error) {
	command := `SELECT id,
					  value,
					  additional_value
					FROM client_benefit_user
					WHERE client_id = $1 AND client_benefit_id = $2 `
	err = br.Conn.Select(&benefitsUser, command, bu.ClientID, bu.BenefitID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return benefitsUser, err
}

// def get_users(db: Session, skip: int = 0, limit: int = 100):
//     return db.query(models.User).offset(skip).limit(limit).all()

// def create_user(db: Session, bs_user: schemas.UserCreate):
//     fake_hashed_password = bs_user.password + "notreallyhashed"
//     db_user = models.User(email=bs_user.email, hashed_password=fake_hashed_password)
//     db.add(db_user)
//     db.commit()
//     db.refresh(db_user)
//     return db_user

// def get_items(db: Session, skip: int = 0, limit: int = 100):
//     return db.query(models.Item).offset(skip).limit(limit).all()

// def create_user_item(db: Session, bs_item: schemas.ItemCreate, userID int):
//     db_item = models.Item(**bs_item.dict(), owner_id=user_id)
//     db.add(db_item)
//     db.commit()
//     db.refresh(db_item)
//     return db_item
