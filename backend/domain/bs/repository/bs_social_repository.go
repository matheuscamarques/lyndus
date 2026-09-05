package repository

import (
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/infra/db"
	"database/sql"
)

type BsSocialRepository struct {
	db.Connector
}

func NewBSSocialRepository() (repo BsSocialRepository) {
	repo.Table = "bs_social"
	repo.UsePostgres()
	return repo
}

func (b BsSocialRepository) GetSocials(bsID int) (socials []entity.Social, err error) {
	command := `SELECT soc.id,
       				   soc."desc"
					FROM bs_social soc
    				WHERE soc.id NOT IN (SELECT bs.social_id 
    										FROM bs_social bs
    										WHERE bs_id = $1)`
	err = b.Conn.Select(&socials, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return socials, err
}

func (b BsSocialRepository) GetSocialByID(socialID int) (social entity.BsSocial, err error) {
	command := `SELECT soc.id as social_id
					FROM bs_social soc 
					WHERE id = $1`
	err = b.Conn.Get(&social, command, socialID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return social, err
}

func (b BsSocialRepository) GetBSSocialByID(bsID, socialID int) (social entity.BsSocial, err error) {
	command := `SELECT bss.id
					FROM bs_social bss 
					WHERE bss.social_id = $1
					AND bss.bs_id = $2`
	err = b.Conn.Get(&social, command, socialID, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return social, err
}

func (b BsSocialRepository) GetAll(bsID int) (socials []entity.BsSocial, err error) {
	command := `SELECT bss.bs_id,
       				   bss.social_id,
       				   bss.url,
       				   soc."desc"
       				FROM bs_social bss
					INNER JOIN bs_social soc on bss.social_id = soc.id
					WHERE bs_id = $1`
	err = b.Conn.Select(&socials, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return socials, err
}

func (b BsSocialRepository) Update(bsSocial entity.BsSocial) error {
	command := `UPDATE bs_social SET url = $1 
					WHERE bs_id = $2
					AND social_id = $3`
	stmt, err := b.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		bsSocial.Url,
		bsSocial.BSID,
		bsSocial.SocialID,
	)
	return err
}

func (b BsSocialRepository) Create(bsSocial entity.BsSocial) (id int, err error) {
	command := `INSERT INTO bs_social(social_id, bs_id, url) 
				VALUES ($1,$2,$3)
                RETURNING ID`
	stmt, err := b.Conn.Prepare(command)
	if err != nil {
		return id, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(bsSocial.SocialID, bsSocial.BSID, bsSocial.Url).Scan(&id)
	return id, err
}

func (b BsSocialRepository) Remove(bsID int, socialID int) error {
	command := `DELETE FROM bs_social WHERE bs_id = $1 AND social_id = $2`
	stmt, err := b.Conn.Prepare(command)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		bsID,
		socialID,
	)
	return err

}
