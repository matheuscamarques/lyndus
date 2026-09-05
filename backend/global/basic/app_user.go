package basic

import "bitbucket.org/lyndus/backend/infra/types"

type AppUser struct {
	ID       int       `db:"id" json:"id"`
	PersonID int       `db:"person_id" json:"personID"`
	Cpf      *types.CPF `db:"cpf" json:"cpf"`
	Email    *string    `db:"email" json:"email"`
}
