package entity

import "bitbucket.org/lyndus/backend/infra/utils"

type ClientUser struct {
	ID               int     `json:"id" db:"id" `
	Phone            *string `json:"phone,omitempty" db:"phone"`
	Username         string  `json:"username" db:"username"`
	Name             string  `json:"name"   db:"name"`
	Email            string  `json:"email"  db:"email"`
	StatusID         int     `json:"-" db:"status"`
	Status           string  `json:"status" db:"status"`
	Active           bool    `json:"active"`
	ClientID         int     `json:"-" db:"client_id"`
	CreatedAt        int     `json:"-" db:"created_at"`
	CreatedBy        int     `json:"-" db:"created_by"`
	Master           bool    `json:"-" db:"master"`
	AuthenticationID int     `json:"-" db:"authentication_id"`
	Password         string  `json:"password,omitempty"db:"password"`

	Permissions []utils.LevelPermission `json:"permissions,omitempty"`
}
