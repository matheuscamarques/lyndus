package entity

type BSBanner struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"bane" db:"name"`
	URL  string `json:"url" db:"url"`
}
