package repository

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type AppUserBSBannerRepository struct {
	conn *sqlx.DB
}

func NewAppUserBSBannerRepository() AppUserBSBannerRepository {
	return AppUserBSBannerRepository{
		conn: postgres.DB,
	}
}


func (w AppUserBSBannerRepository) GetBSBanners(bsID int) (banners []entity.BSBanner, err error) {
	command := `SELECT id, name FROM bs_banners WHERE deleted=false AND bs_id=$1`
	err = w.conn.Select(&banners, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}