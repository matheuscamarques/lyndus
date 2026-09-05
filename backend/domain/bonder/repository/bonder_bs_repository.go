package repository

import (
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/infra/db"
	"bitbucket.org/lyndus/backend/infra/types"
	"bitbucket.org/lyndus/backend/internal/constants"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"

	"github.com/jmoiron/sqlx"
	"log"
)

type BonderBsRepository struct {
	db.Connector
}

func (bbr BonderBsRepository) GetConnection() *sqlx.DB {
	return bbr.Conn
}

func (bbr BonderBsRepository) GetConnector() db.Connector {
	return bbr.Connector
}

func NewBonderBsRepository() (repo BonderBsRepository) {
	repo.UsePostgres()
	repo.Table = "bs"
	return repo
}

func (bbr BonderBsRepository) GetByIdSimple(id int) (bs entity.BS, err error) {
	command := `SELECT bs.id,
       				   bs.company_id,
       				   bs.cnpj
					FROM bs 
					WHERE bs.id = $1
					`
	err = bbr.Conn.Get(&bs, command, id)
	return bs, err
}

func (bbr BonderBsRepository) GetById(id int) (bs entity.BS, err error) {
	command := `SELECT 
						bs.id,
       					bs.lat,
       					bs.lon,
						com.company_name,
						com.fantasy_name,
						com.cnpj,
       					com.phone,
       					com.email,
       					com.state,
       					com.city,
       					com.district,
       					com.street,
       					com.number,
       					com.address_complement,
       					com.zipcode,
       					com."desc",
       					bs.active_lyndus,
					    bs."active",
       					bs.plan_value,
       					bs.plan_day
					FROM bs 
					INNER JOIN company com ON bs.company_id = com.id
					WHERE bs.id = $1
					`
	err = bbr.Conn.QueryRow(command, id).Scan(
		&bs.ID,
		&bs.Lat,
		&bs.Lon,
		&bs.CompanyName,
		&bs.FantasyName,
		&bs.CNPJ,
		&bs.Phone,
		&bs.Email,
		&bs.State,
		&bs.City,
		&bs.District,
		&bs.Street,
		&bs.Number,
		&bs.AddressComplement,
		&bs.Zipcode,
		&bs.Desc,
		&bs.Active,
		&bs.ActiveApp,
		&bs.PlanValue,
		&bs.PlanDay)

	return bs, err
}

func (bbr BonderBsRepository) GetLyndusPay(bsID int) (closingDay int, rateAnticipation, rateAnticipationBm decimal.Decimal, err error) {

	command := `SELECT bar.closing_day,
					   CASE WHEN bar.rate_anticipation is null then 0 else bar.rate_anticipation END,
					   CASE WHEN bar.rate_anticipation_bm is null then 0 else bar.rate_anticipation_bm END 
					FROM bs_account_receipt bar
					    INNER JOIN bs_account_category bac on bar.bs_account_category_id = bac.id
						LEFT JOIN bs_account_movement bam on bar.bs_account_movement_id = bam.id
						WHERE bar.bs_id = $1
						  AND bar.active = true
						  AND bar.deleted = false
						  AND bac.id = $2
					ORDER BY bar.name `

	err = bbr.Conn.QueryRow(command, bsID, constants.RecebimentoLyndusPay).Scan(
		&closingDay,
		&rateAnticipation,
		&rateAnticipationBm)
	if err == pgx.ErrNoRows {
		return closingDay, rateAnticipation, rateAnticipationBm, nil
	} else if err != nil {
		return closingDay, rateAnticipation, rateAnticipationBm, err
	}

	return closingDay, rateAnticipation, rateAnticipationBm, nil
}
func (bbr BonderBsRepository) GetAll(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (bss []entity.BasicBS, totalItems, totalPages int, err error) {
	command := `SELECT count(*)
					FROM bs 
					INNER JOIN company com ON com.id = bs.company_id
					WHERE bs.deleted=false
					  AND bs.active_lyndus=$1 `

	if search != "" {
		command = command + ` AND (unaccent(com.company_name) ILIKE $2 
							   OR unaccent(com.fantasy_name) ILIKE $2) `
		err = bbr.Conn.Get(&totalItems, command, active, "%"+search+"%")
	} else {
		err = bbr.Conn.Get(&totalItems, command, active)
	}

	if err != nil {
		return bss, totalItems, totalPages, err
	}
	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT 
					bs.id,
					com.company_name,
					com.fantasy_name,
					com.cnpj
					FROM bs 
					INNER JOIN company  com ON bs.company_id = com.id
					WHERE bs.deleted = false
					  AND bs.active_lyndus = $3 `

	if search != "" {
		command = command + ` AND (unaccent(com.company_name) ILIKE $4 
							   OR unaccent(com.fantasy_name) ILIKE $4) `
	}

	if orderBy == "companyName" {
		command = command + ` ORDER BY com.company_name `
	} else if orderBy == "fantasyName" {
		command = command + ` ORDER BY com.fantasy_name `
	} else if orderBy == "cnpj" {
		command = command + ` ORDER BY com.cnpj `
	} else {
		command = command + ` ORDER BY com.company_name `
	}

	if sortDesc {
		command = command + " DESC "
	}
	command = command + ` OFFSET $1 LIMIT $2`

	if search != "" {
		err = bbr.Conn.Select(&bss, command, (activePage-1)*itemsPerPage, itemsPerPage, active, "%"+search+"%")
	} else {
		err = bbr.Conn.Select(&bss, command, (activePage-1)*itemsPerPage, itemsPerPage, active)
	}
	if err == sql.ErrNoRows {
		err = nil
	}

	return bss, totalItems, totalPages, err
}

func (bbr BonderBsRepository) GetTable() string {
	return bbr.Table
}

func (bbr BonderBsRepository) Create(bs entity.BS) (int, error) {
	command := "INSERT INTO bs(company_id,cnpj,lat,lon, plan_value, plan_day, balance_day, rate_anticipation, rate_anticipation_bm) VALUES($1,$2,$3,$4,$5,$6,$7,$8, $9) RETURNING id"
	stmt, err := bbr.Conn.Prepare(command)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var id int

	err = stmt.QueryRow(bs.CompanyID, bs.CNPJ, bs.Lat, bs.Lon, bs.PlanValue, bs.PlanDay, bs.BalanceDay, bs.RateAnticipation, bs.RateAnticipationBm).Scan(&id)
	if err != nil {
		return 0, err
	}

	err = bbr.GeoNotifyUpdate(id)
	return id, err
}

func (bbr BonderBsRepository) CreateBasicAccounts(bsID, closingDay int, rateAnticipation decimal.Decimal) error {
	command := `insert into bs_account_movement ( bs_id, bs_account_category, name, bank, agency, bank_account, initial_balance, balance_date, chave_pix, chave_pix_type ) values
		                                ( $1, 3, 'CAIXA', '', '', '', 0, now(), '', '')`
	_, err := bbr.Conn.Exec(command, bsID)
	if err != nil {
		return err
	}

	command = `insert into bs_account_receipt (bs_id, bs_account_category_id, name, fee, card_machine, value_machine, day_machine, installment_machine, installment_fee, card_cycle, rescue_day, closing_day, receipt_day, installment, installment_max, installment_fee_by, free_fee, installments_free_fee, rate_anticipation) values
                               ($1, 6, 'Aplicativo', 0.00, '', 0.00, 0, 0, null, '', 0, $2, 6, null, null, null, null, null, $3);`
	_, err = bbr.Conn.Exec(command, bsID, closingDay, rateAnticipation)

	if err != nil {
		return err
	}

	command = `insert into bs_account_receipt (bs_id, bs_account_category_id, name, fee, card_machine, value_machine, day_machine, installment_machine, installment_fee, card_cycle, rescue_day, closing_day, receipt_day, installment, installment_max, installment_fee_by, free_fee, installments_free_fee, rate_anticipation) values
                               ($1, 7, 'Maquininha', 0.00, '', 0.00, 0, 0, null, '', 0, $2, 6, null, null, null, null, null, $3);`
	_, err = bbr.Conn.Exec(command, bsID, closingDay, rateAnticipation)

	return err
}
func (bbr BonderBsRepository) GeoNotifyUpdate(bsID int) (err error) {
	command := `UPDATE bs SET 
              				geog=ST_SetSRID(ST_MakePoint(bs.lon, bs.lat), 4326),
              				point=ST_MakePoint(bs.lon, bs.lat) 
			   WHERE id = $1`
	_, err = bbr.Conn.Exec(command, bsID)
	return
}

func (bbr BonderBsRepository) SearchByCNPJ(cnpj types.CNPJ) (id int, err error) {
	command := "SELECT id FROM bs WHERE cnpj = $1 "
	err = bbr.Conn.Get(&id, command, cnpj)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return 0, err
	}
	return id, err
}

func (bbr BonderBsRepository) Update(bs entity.BS) (err error) {
	command := `
		UPDATE bs_account_receipt SET	
			closing_day=$4,
		    rate_anticipation=$5,
		    rate_anticipation_bm=$6
		WHERE bs_id = $1 
		AND (bs_account_category_id = $2 OR bs_account_category_id = $3)
	`
	_, err = bbr.Conn.Exec(command, bs.ID, constants.RecebimentoLyndusPay, constants.RecebimentoLyndusPayCartão, bs.BalanceDay, bs.RateAnticipation, bs.RateAnticipationBm)
	if err != nil {
		return err
	}

	command = `
		UPDATE bs SET	
			lat=$2,
		    lon=$3,
		    plan_value=$4,
		    plan_day=$5
		WHERE id = $1		
	`
	_, err = bbr.Conn.Exec(command, bs.ID, bs.Lat, bs.Lon, bs.PlanValue, bs.PlanDay)
	if err != nil {
		return err
	}
	return bbr.GeoNotifyUpdate(bs.ID)
}

func (bbr BonderBsRepository) UpdateLyndusActive(bsId int, active bool) (err error) {
	command := `
		UPDATE bs SET	
			active_lyndus=$2
		WHERE id = $1
	`
	log.Println(command, bsId, active)
	_, err = bbr.Conn.Exec(command, bsId, active)

	return err
}

func (bbr BonderBsRepository) Detete(bsId int) (err error) {
	command := `
		UPDATE bs SET	
			deleted=$2
		WHERE id = $1		
	`
	_, err = bbr.Conn.Exec(command, bsId, true)

	return err
}

func (bbr BonderBsRepository) GetFinancialList(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (wcsList []entity.WithdrawalCashSimple, totalItems, totalPages int, err error) {
	command := `SELECT count(*)
					FROM bs_withdrawal_cash bwc
					INNER JOIN bs ON bwc.bs_id = bs.id
					INNER JOIN company com ON bs.company_id = com.id `

	if search != "" {
		command = command + ` AND unaccent(com.company_name) ILIKE unaccent($1) `
		err = bbr.Conn.QueryRow(command, "%"+search+"%").Scan(&totalItems)
	} else {
		err = bbr.Conn.QueryRow(command).Scan(&totalItems)
	}
	if err != nil {
		return wcsList, totalItems, totalPages, err
	}

	rest := totalItems % itemsPerPage
	totalPages = totalItems / itemsPerPage
	if rest > 0 {
		totalPages += 1
	}
	command = `SELECT bwc.id,
					  bwc.bs_id,
					  bwc.value_released,
					  bwc.value_anticipation,
					  bwc.total_value,
					  bwc.request_date,
					  bwc.pay_date,
					  bwc.requested,
					  bwc.paid,
       				  com.company_name,
					  bs.cnpj
					FROM bs_withdrawal_cash bwc
					INNER JOIN bs ON bwc.bs_id = bs.id
					INNER JOIN company com ON bs.company_id = com.id `

	if search != "" {
		command = command + ` AND unaccent(com.company_name) ILIKE unaccent($3) `
	}

	if orderBy == "payDate" {
		command = command + ` ORDER BY bwc.pay_date `
	} else if orderBy == "valueReseased" {
		command = command + ` ORDER BY bwc.value_released `
	} else if orderBy == "valueAnticipation" {
		command = command + ` ORDER BY bwc.value_anticipation `
	} else if orderBy == "totalValue" {
		command = command + ` ORDER BY bwc.total_value `
	} else if orderBy == "requestDate" {
		command = command + ` ORDER BY bwc.request_date `
	} else if orderBy == "paid" {
		command = command + ` ORDER BY bwc.paid `
	} else {
		command = command + ` ORDER BY bwc.requested ASC, bwc.request_date `
	}

	if sortDesc {
		command = command + " DESC "
	}
	command = command + ` OFFSET $1 LIMIT $2`

	var rows *sql.Rows

	if search != "" {
		rows, err = bbr.Conn.Query(command, (activePage-1)*itemsPerPage, itemsPerPage, "%"+search+"%")
	} else {
		rows, err = bbr.Conn.Query(command, (activePage-1)*itemsPerPage, itemsPerPage)
	}

	defer rows.Close()
	if err == pgx.ErrNoRows {
		return wcsList, totalItems, totalPages, err
	}
	if err != nil {
		return wcsList, totalItems, totalPages, err
	}
	var wcs entity.WithdrawalCashSimple
	for rows.Next() {
		wcs = entity.WithdrawalCashSimple{}
		err = rows.Scan(
			&wcs.ID,
			&wcs.BSID,
			&wcs.ValueReleased,
			&wcs.ValueAnticipation,
			&wcs.TotalValue,
			&wcs.RequestDate,
			&wcs.PayDate,
			&wcs.Requested,
			&wcs.Paid,
			&wcs.CompanyName,
			&wcs.CNPJ)
		if err != nil {
			return wcsList, totalItems, totalPages, err
		}
		wcsList = append(wcsList, wcs)
	}
	err = rows.Err()

	return wcsList, totalItems, totalPages, err
}

func (bbr BonderBsRepository) GetBSWithdrawalCash(id int) (withdrawalCash entity.WithdrawalCash, err error) {
	command := `SELECT bwc.id,
					   bwc.bs_id,
					   bwc.bs_bank_details_id,
					   bwc.created_at,
					   bwc.value_released,
					   bwc.value_anticipation,
					   bwc.rate_anticipation,
					   bwc.total_value,
					   bwc.request_date,
					   bwc.pay_date,
					   bwc.requested,
					   bwc.paid,
					   bwc.balance_day,
					   bam.name,
					   bam.bank,
					   bam.agency,
					   bam.bank_account,
					   bam.chave_pix,
       				   com.company_name,
					   bs.cnpj,
					   bwc.bs_account_movement_id,
					   bwc.bs_account_receipt_id
					FROM bs_withdrawal_cash bwc 
						INNER JOIN bs_account_movement bam ON bwc.bs_account_movement_id = bam.id
						INNER JOIN bs ON bwc.bs_id = bs.id
						INNER JOIN company com ON bs.company_id = com.id 
						WHERE  bwc.id = $1 `

	err = bbr.Conn.QueryRow(command, id).Scan(
		&withdrawalCash.ID,
		&withdrawalCash.BSID,
		&withdrawalCash.BankDetailsID,
		&withdrawalCash.CreatedAt,
		&withdrawalCash.ValueReleased,
		&withdrawalCash.ValueAnticipation,
		&withdrawalCash.RateAnticipation,
		&withdrawalCash.TotalValue,
		&withdrawalCash.RequestDate,
		&withdrawalCash.PayDate,
		&withdrawalCash.Requested,
		&withdrawalCash.Paid,
		&withdrawalCash.BalanceDay,
		&withdrawalCash.AccountName,
		&withdrawalCash.Bank,
		&withdrawalCash.Agency,
		&withdrawalCash.BankAccount,
		&withdrawalCash.ChavePIX,
		&withdrawalCash.CompanyName,
		&withdrawalCash.CNPJ,
		&withdrawalCash.AccountMovementID,
		&withdrawalCash.AccountReceiptID)

	if errors.Is(err, sql.ErrNoRows) {
		command = `SELECT bwc.id,
					   bwc.bs_id,
					   bwc.bs_bank_details_id,
					   bwc.created_at,
					   bwc.value_released,
					   bwc.value_anticipation,
					   bwc.rate_anticipation,
					   bwc.total_value,
					   bwc.request_date,
					   bwc.pay_date,
					   bwc.requested,
					   bwc.paid,
					   bwc.balance_day,
					   bbd.bank,
					   bbd.agency,
					   bbd.bank_account,
					   bbd.chave_pix,
       				   com.company_name,
					   bs.cnpj,
					   bwc.bs_account_movement_id,
					   bwc.bs_account_receipt_id
					FROM bs_withdrawal_cash bwc 
						INNER JOIN bs_bank_details bbd ON bwc.bs_bank_details_id = bbd.id
						INNER JOIN bs ON bwc.bs_id = bs.id
						INNER JOIN company com ON bs.company_id = com.id 
						WHERE  bwc.id = $1 `

		err = bbr.Conn.QueryRow(command, id).Scan(
			&withdrawalCash.ID,
			&withdrawalCash.BSID,
			&withdrawalCash.BankDetailsID,
			&withdrawalCash.CreatedAt,
			&withdrawalCash.ValueReleased,
			&withdrawalCash.ValueAnticipation,
			&withdrawalCash.RateAnticipation,
			&withdrawalCash.TotalValue,
			&withdrawalCash.RequestDate,
			&withdrawalCash.PayDate,
			&withdrawalCash.Requested,
			&withdrawalCash.Paid,
			&withdrawalCash.BalanceDay,
			&withdrawalCash.Bank,
			&withdrawalCash.Agency,
			&withdrawalCash.BankAccount,
			&withdrawalCash.ChavePIX,
			&withdrawalCash.CompanyName,
			&withdrawalCash.CNPJ,
			&withdrawalCash.AccountMovementID,
			&withdrawalCash.AccountReceiptID)
	}

	return withdrawalCash, err
}

func (bbr BonderBsRepository) GetBSWithdrawalCashBalance(bsID, id int) (wcbs []entity.WithdrawalCashBalance, err error) {
	command := `SELECT bs_balance_receivable_id,
       				   bs_withdrawal_cash_id,
       				   value
					FROM bs_withdrawal_cash_balance_receivable bwcbr 
					INNER JOIN bs_withdrawal_cash bwc ON bwc.id = bwcbr.bs_withdrawal_cash_id 
					WHERE bwc.bs_id = $1
					  AND bwc.id = $2`

	rows, err := bbr.Conn.Query(command, bsID, id)
	if err != nil {
		return wcbs, err
	}
	var wcb entity.WithdrawalCashBalance
	for rows.Next() {
		wcb = entity.WithdrawalCashBalance{}
		err = rows.Scan(
			&wcb.BalanceReceivableID,
			&wcb.WithdrawalCashID,
			&wcb.Value)
		wcbs = append(wcbs, wcb)
	}

	return wcbs, err
}

func (bbr BonderBsRepository) BSWithdrawalCashConfirm(bsID int, accountMovementID *int, accountReceiptID int, value, valueRate decimal.Decimal, wcbs []entity.WithdrawalCashBalance) (err error) {
	command := `UPDATE bs_balance_receivable SET 
							bs_paid = true
						WHERE bs_paid_value = value
						  AND id = $1
						  AND bs_id = $2 `

	tx, err := bbr.Conn.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	for k := range wcbs {
		_, err = tx.Exec(command, wcbs[k].BalanceReceivableID, bsID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	command = `UPDATE bs_withdrawal_cash SET 
							paid = true, 
							pay_date = now()
						WHERE id = $1
						  AND bs_id = $2 `

	_, err = tx.Exec(command, wcbs[0].WithdrawalCashID, bsID)
	if err != nil {
		tx.Rollback()
		return err
	}

	insertCashFlow := `INSERT INTO bs_cash_flow (bs_id,
                          						 competency_date,
                          						 due_date,
                          						 movement_date,
                          						 description,
                          						 bs_account_movement_id,
                          						 bs_account_receipt_id,
                          						 finance_category_id,
                          						 "value",
                          						 movement_type
                          						 )
						VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	//registrar movimento no caixa - bs_cash_flow
	var today types.Date
	today.SetNow()

	_, err = tx.Exec(insertCashFlow,
		bsID,
		today,
		today,
		today,
		"Resgate Lyndus Bank",
		accountMovementID,
		accountReceiptID,
		constants.FinanceRecebimentoLyndusPay,
		value.Sub(valueRate),
		constants.MovementTypeDebito,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if valueRate.IsPositive() {
		_, err = tx.Exec(insertCashFlow,
			bsID,
			today,
			today,
			today,
			"Taxa de Antecipação LyndusPay",
			accountMovementID,
			accountReceiptID,
			constants.FinanceTaxaAntecipacaoLyndusPay,
			valueRate,
			constants.MovementTypeDebito,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}
