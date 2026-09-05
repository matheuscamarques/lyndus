package entity

type ClientPollCategory struct {
	ID int `json:"id"`
	ClientPollID int `json:"client_poll_id"`
	ClientCategoryID int `json:"client_category_id"`
}
