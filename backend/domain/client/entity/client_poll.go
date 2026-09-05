package entity

import "bitbucket.org/lyndus/backend/infra/types"

type ClientPoll struct {
	ID             int            `db:"id"`
	ClientID       int            `db:"client_id"`
	StartTime      types.DateTime `db:"start_time"`
	EndTime        types.DateTime `db:"end_time"`
	AllEmployee    bool           `db:"all_employee"`
	Title          string         `db:"title"`
	Status         uint8            `db:"status"`
	CreatedBy	   int			  `db:"created_by"`
}
