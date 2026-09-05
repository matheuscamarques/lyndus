package entity

type BSBankDetails struct {
	ID          int    `json:"_" db:"id"`
	BsID        int    `json:"_" db:"bs_id"`
	Bank        string `json:"bank" db:"bank"`
	Branch      string `json:"branch" db:"branch"`
	BankAccount string `json:"bankAccount" db:"bank_account"`
}
