package entity

type Authentication struct {
	ID int `json:"id"`
	Password string `json:"password"`
	Disable bool `json:"disable"`
}