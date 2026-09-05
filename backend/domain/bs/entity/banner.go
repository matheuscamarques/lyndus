package entity

type Banner struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Image string `json:"image" db:"image"`
}