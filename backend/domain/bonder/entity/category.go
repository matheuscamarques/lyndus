package entity

type Category struct {
	ID   int    `json:"id" db:"id"`
	Desc string `json:"desc" db:"desc"`
}
