package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type BSBannerRepository struct {
	contracts.BSBannerRepositoryInterface
	conn  *sqlx.DB
}

func  NewBSBannerRepository() BSBannerRepository {
	return BSBannerRepository{
		conn : postgres.DB,
	}
}

//CreateBanner func to create banner
func (b BSBannerRepository) CreateBanner(bsID int, name string) (id int, err error) {
	command := `INSERT INTO bs_banners(bs_id, "name") VALUES($1, $2) RETURNING ID`
	stmt, err := b.conn.Prepare(command)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(bsID, name).Scan(&id)
	return id, err
}

func (b BSBannerRepository) GetBanners(bsID int) (banners []entity.Banner, err error) {
	command := `SELECT id, "name" FROM bs_banners WHERE bs_id=$1 AND deleted = false`
	err = b.conn.Select(&banners, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return banners, err
}

func (b BSBannerRepository) DeleteBanner(bsID, bannerID int) error {
	command := `UPDATE bs_banners SET deleted = true WHERE id=$1 AND bs_id=$2 `
	stmt, err := b.conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(bannerID, bsID)
	return err
}