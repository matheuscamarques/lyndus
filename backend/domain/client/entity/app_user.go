package entity

import (
	"bitbucket.org/lyndus/backend/infra/types"
	"time"
)

type AppUser struct {
	ID               int        `db:"id" json:"id" `
	Name             string     `db:"name" json:"name"`
	PersonID         int        `db:"person_id" json:"personID"`
	CPF              types.CPF  `db:"cpf" json:"cpf"`
	Email            string     `db:"email" json:"email"`
	CreatedAt        time.Time  `db:"created_at,omitempty" json:"created_at,omitempty"`
	AuthenticationID int        `db:"authentication_id" json:"authentication_id,omitempty"`
	AcceptedTerm     bool       `db:"accepted_term" json:"accepted_term,omitempty"`
}