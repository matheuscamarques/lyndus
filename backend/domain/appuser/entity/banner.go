package entity

import (
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
)

type Banner struct {
	ID          int     `json:"id" db:"id"`
	Image       string  `json:"image" db:"image"`
	MobileImage *string `json:"mobileImage" db:"mobile_image"`
	Description *string `json:"description" db:"description"`
}

func (b Banner) GetAppBannerImages() (bss []Banner, err error) {
	command := `SELECT id, image, mobile_image, description FROM app_banner`
	err = postgres.DB.Select(&bss, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bss, err
}
