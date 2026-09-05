package entity

import "bitbucket.org/lyndus/backend/infra/utils"

type BSUser struct {
	ID               int     `json:"id" db:"id" `
	Phone            *string `json:"phone" db:"phone"`
	Username         string  `json:"username" db:"username"`
	Name             string  `json:"name"   db:"name"`
	Email            string  `json:"email"  db:"email"`
	Status           int     `json:"status,omitempty" db:"status"`
	BsID             int     `json:"-" db:"bs_id"`
	CreatedAt        int     `json:"-" db:"created_at"`
	CreatedBy        int     `json:"-" db:"created_by"`
	Master           bool    `json:"-" db:"master"`
	AuthenticationID int     `json:"-" db:"authentication_id"`
	Password         string  `json:"password,omitempty"`

	Permissions []utils.LevelPermission `json:"permissions,omitempty" db:"password"`
}
