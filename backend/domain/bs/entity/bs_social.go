package entity

type Social struct {
	ID   int    `json:"id" db:"id"`
	Desc string `json:"desc" db:"desc"`
}

type BsSocial struct {
	ID       int    `json:"id,omitempty" db:"id"`
	SocialID int    `json:"socialID" db:"social_id"`
	BSID     int    `json:"-" db:"bs_id"`
	Url      string `json:"url" db:"url"`
	Desc     string `json:"desc" db:"desc"`
}
