package repository

import (
	entity2 "bitbucket.org/lyndus/backend/domain/bonder/entity"
	"database/sql"
	"time"

	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/db"

	"bitbucket.org/lyndus/backend/domain/appuser/entity"
)

const SRID = 4326

type AppUserBSRepository struct {
	db.Connector
}

func NewAppUserBSRepository() (repo *AppUserBSRepository) {
	repo = &AppUserBSRepository{}
	repo.Table = "bs"
	repo.UsePostgres()
	return repo
}

func (aubr AppUserBSRepository) GetLocation(bsID int) (location string, err error) {
	command := `SELECT location
					FROM bs 
					WHERE bs.id = $1 `
	err = aubr.Conn.Get(&location, command, bsID)
	return location, err
}

func (aubr AppUserBSRepository) GetBSName(bsId int) (name string, err error) {
	command := `SELECT com.fantasy_name
                FROM bs
                INNER JOIN company com ON bs.company_id = com.id
                WHERE BS.id = $1`
	err = aubr.Conn.Get(&name, command, bsId)
	return name, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSBMName(bsId, bmID int) (name string, err error) {
	command := `SELECT bsb.name
                FROM bs
                INNER JOIN bs_bm bsb ON bs.id=bsb.bs_id
                WHERE bs.id = $1 AND bsb.id = $2 `
	err = aubr.Conn.Get(&name, command, bsId, bmID)
	return name, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBS(bsId int) (bs entity.BS, err error) {
	command := `SELECT bs.id,
					   com.cnpj,
					   com.company_name,
					   com.phone,
					   com.fantasy_name,
					   com.email,
					   com.state,
					   com.city, 
					   com.district,
					   com.street,
					   com.number,
					   com.address_complement,
					   com.zipcode,
					   com.desc,
					   bs.lat,
					   bs.lon
                FROM bs
                INNER JOIN company com ON bs.company_id = com.id
                WHERE BS.id = $1`
	err = aubr.Conn.Get(&bs, command, bsId)

	command = `SELECT bss.social_id,
       				  soc."desc",
       				  bss.url
				FROM bs_social soc
				INNER JOIN bs_social bss ON soc.id = bss.social_id
				WHERE bss.bs_id=$1`
	err = aubr.Conn.Select(&bs.Socials, command, bsId)

	return bs, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSID(bsId int) (id int, err error) {
	command := `SELECT bs.id
					FROM bs
					WHERE BS.id = $1`
	err = aubr.Conn.Get(&id, command, bsId)
	return id, aubr.Error(err)
}

//todo criar coluna no banco para otimizar busca por fantasy_name. A coluna vai limpar acentos e padronizar as letras sem acento.

func (aubr AppUserBSRepository) GetBSs(fantasyName string, lat, lon float64, criteria criteria.Criteria) (bss []aggregate.BSAggregate, totalPages int, err error) {
	criteria.CheckPages()

	command := `SELECT 
                       count(*)
					FROM bs
					INNER JOIN company com ON bs.company_id = com.id
					INNER JOIN bs_bs_category bbc ON bs.id = bbc.bs_id `

	if lat != 0 && lon != 0 && fantasyName != "" {
		command = command + ` WHERE lower(unaccent(com.fantasy_name)) LIKE lower(unaccent($1)) 
									AND ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($2, $3), 4326), 10000) `
		err = aubr.Conn.Get(&totalPages, command, "%"+fantasyName+"%", lon, lat)
	} else if lat != 0 && lon != 0 {
		command = command + ` WHERE ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($1, $2), 4326), 10000) `
		err = aubr.Conn.Get(&totalPages, command, lon, lat)
	} else if fantasyName != "" {
		command = command + ` WHERE lower(unaccent(com.fantasy_name)) LIKE lower(unaccent($1)) `
		err = aubr.Conn.Get(&totalPages, command, "%"+fantasyName+"%")
	} else {
		err = aubr.Conn.Get(&totalPages, command)
	}

	if err != nil {
		return bss, totalPages, aubr.Error(err)
	}

	rest := totalPages % criteria.ItemsPerPage
	totalPages = totalPages / criteria.ItemsPerPage

	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT 
                       bs.id,
					   com.cnpj,
					   com.company_name "company.company_name",
					   com.phone  		"company.phone",
					   com.fantasy_name "company.fantasy_name",
					   com.email		"company.email",
					   com.state		"company.state",
					   com.city			"company.city",
					   com.district		"company.district",
					   com.street		"company.street",
					   com.number		"company.number",
					   com.address_complement "company.address_complement",
					   com.zipcode			  "company.zipcode",
					   com.desc				  "company.desc" ,
					   bs.lat				  "bs.lat",
					   bs.lon				  "bs.lon"
					FROM bs
					INNER JOIN company com ON bs.company_id = com.id
					INNER JOIN bs_bs_category bbc ON bs.id = bbc.bs_id`

	if lat != 0 && lon != 0 && fantasyName != "" {
		command = command + ` WHERE lower(unaccent(com.fantasy_name)) LIKE lower(unaccent($1)) 
  							  AND ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($2, $3), 4326), 10000) 
  							  ORDER BY ORDER BY bs.geog <-> ST_SetSRID(ST_MakePoint($2, $3), 4326)`
		command = criteria.Query(command)
		err = aubr.Conn.Select(&bss, command, "%"+fantasyName+"%", lon, lat)
	} else if lat != 0 && lon != 0 {
		command = command + ` WHERE ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($1, $2), 4326), 10000) 
							  ORDER BY bs.geog <-> ST_SetSRID(ST_MakePoint($1, $2), 4326) `
		command = criteria.Query(command)
		err = aubr.Conn.Select(&bss, command, lon, lat)
	} else if fantasyName != "" {
		command = command + ` WHERE lower(unaccent(com.fantasy_name)) LIKE lower(unaccent($1)) 
							  ORDER BY com.fantasy_name `
		command = criteria.Query(command)
		err = aubr.Conn.Select(&bss, command, "%"+fantasyName+"%")
	} else {
		command = command + ` ORDER BY com.fantasy_name `
		command = criteria.Query(command)
		err = aubr.Conn.Select(&bss, command)
	}
	if err == sql.ErrNoRows {
		err = nil
		totalPages = 0
	}
	return bss, totalPages, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSsMap(lat, lon float64) (bss []entity.BS, err error) {
	command := `SELECT bs.id,
					   com.cnpj,
					   com.company_name,
					   com.phone,
					   com.fantasy_name,
					   com.email,
					   com.state,
					   com.city, 
					   com.district,
					   com.street,
					   com.number,
					   com.address_complement,
					   com.zipcode,
					   com.desc,
					   bs.lat,
					   bs.lon
                FROM bs
                INNER JOIN company com ON bs.company_id = com.id
                WHERE ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($1, $2), 4326), 10000)
				ORDER BY bs.geog <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)`
	err = aubr.Conn.Select(&bss, command, lon, lat)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bss, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSsByCategory(categoryID int, fantasyName string, lat, lon float64, criteria criteria.Criteria) (bss []aggregate.BSAggregate, totalPages int, err error) {
	criteria.CheckPages()

	command := `SELECT 
                       count(*)
					FROM bs
					INNER JOIN company com ON bs.company_id = com.id
					INNER JOIN bs_bs_category bbc ON bs.id = bbc.bs_id
					WHERE bbc.bs_category_id=$1 `

	if lat != 0 && lon != 0 && fantasyName != "" {
		command = command + ` AND lower(unaccent(com.fantasy_name)) LIKE lower(unaccent($2)) AND ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($3, $4), 4326), 10000) `
		err = aubr.Conn.Get(&totalPages, command, categoryID, "%"+fantasyName+"%", lon, lat)
	} else if lat != 0 && lon != 0 {
		command = command + ` AND ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($2, $3), 4326), 10000) `
		err = aubr.Conn.Get(&totalPages, command, categoryID, lon, lat)
	} else if fantasyName != "" {
		command = command + ` AND lower(unaccent(com.fantasy_name)) LIKE lower(unaccent($2)) `
		err = aubr.Conn.Get(&totalPages, command, categoryID, "%"+fantasyName+"%")
	} else {
		err = aubr.Conn.Get(&totalPages, command, categoryID)
	}

	if err != nil {
		return bss, totalPages, aubr.Error(err)
	}

	rest := totalPages % criteria.ItemsPerPage
	totalPages = totalPages / criteria.ItemsPerPage

	if rest > 0 {
		totalPages += 1
	}

	command = `SELECT 
                       bs.id,
					   com.cnpj,
					   com.company_name "company.company_name",
					   com.phone  		"company.phone",
					   com.fantasy_name "company.fantasy_name",
					   com.email		"company.email",
					   com.state		"company.state",
					   com.city			"company.city",
					   com.district		"company.district",
					   com.street		"company.street",
					   com.number		"company.number",
					   com.address_complement "company.address_complement",
					   com.zipcode			  "company.zipcode",
					   com.desc				  "company.desc" ,
					   bs.lat				  "bs.lat",
					   bs.lon				  "bs.lon"
					FROM bs
					INNER JOIN company com ON bs.company_id = com.id
					INNER JOIN bs_bs_category bbc ON bs.id = bbc.bs_id
					WHERE bbc.bs_category_id=$1`

	if lat != 0 && lon != 0 && fantasyName != "" {
		command = command + ` AND lower(unaccent(com.fantasy_name)) LIKE lower(unaccent($2)) 
  							  AND ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($3, $4), 4326), 10000) 
  							  ORDER BY ORDER BY bs.geog <-> ST_SetSRID(ST_MakePoint($3, $4), 4326) `
		command = criteria.Query(command)
		err = aubr.Conn.Select(&bss, command, categoryID, "%"+fantasyName+"%", lon, lat)
	} else if lat != 0 && lon != 0 {
		command = command + ` AND ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($2, $3), 4326), 10000) 
							  ORDER BY bs.geog <-> ST_SetSRID(ST_MakePoint($2, $3), 4326) `
		command = criteria.Query(command)
		err = aubr.Conn.Select(&bss, command, categoryID, lon, lat)
	} else if fantasyName != "" {
		command = command + ` AND lower(unaccent(com.fantasy_name)) LIKE lower(unaccent($2))
							  ORDER BY com.fantasy_name `
		command = criteria.Query(command)
		err = aubr.Conn.Select(&bss, command, categoryID, "%"+fantasyName+"%")
	} else {
		command = command + ` ORDER BY com.fantasy_name `
		command = criteria.Query(command)
		err = aubr.Conn.Select(&bss, command, categoryID)
	}
	if err == sql.ErrNoRows {
		err = nil
		totalPages = 0
	}
	return bss, totalPages, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSsByCategoryMap(categoryID int, lat, lon float64) (bss []entity.BS, err error) {
	command := `SELECT bs.id,
					   com.cnpj,
					   com.company_name,
					   com.phone,
					   com.fantasy_name,
					   com.email,
					   com.state,
					   com.city,
					   com.district,
					   com.street,
					   com.number,
					   com.address_complement,
					   com.zipcode,
					   com.desc,
					   bs.lat,
					   bs.lon
					FROM bs
					INNER JOIN company com ON bs.company_id = com.id
					INNER JOIN bs_bs_category bbc ON bs.id = bbc.bs_id
					WHERE bbc.bs_category_id=$1
					AND ST_DWithin(bs.geog, ST_SetSRID(ST_MakePoint($2, $3), 4326), 10000)
					ORDER BY bs.geog <-> ST_SetSRID(ST_MakePoint($2, $3), 4326)`
	err = aubr.Conn.Select(&bss, command, categoryID, lon, lat)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bss, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSWeekDays(bsID int) (weekDays []entity.WeekDayBS, err error) {
	command := `SELECT wed.id,
					   bwd.start_time,
					   bwd.end_time,
       				   bwd.week_day_id
				FROM weekday wed
				INNER JOIN bs_week_day bwd ON wed.id = bwd.week_day_id
				WHERE bwd.bs_id=$1 ORDER BY wed.sequence`
	err = aubr.Conn.Select(&weekDays, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return weekDays, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSCategories(bsID int) (categories []entity2.Category, err error) {
	command := `SELECT bca.id,
       				   bca.desc
				FROM bs_category bca
				INNER JOIN bs_bs_category bbc ON bca.id = bbc.bs_category_id
				WHERE bbc.bs_id=$1 `
	err = aubr.Conn.Select(&categories, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return categories, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSDaysOff(bsID int, start, end time.Time) (days []entity.DayOff, err error) {
	command := `SELECT id, 
       				   start_date,
       				   end_date
                FROM bs_day_off
                WHERE bs_id=$1
                    AND start_date >= $2
  					AND end_date <= $3
                	OR (DATE($2) between DATE(start_date) AND DATE(end_date)
                  		OR DATE($3) BETWEEN DATE(start_date) AND DATE(end_date))`
	err = aubr.Conn.Select(&days, command, bsID, start, end)
	if err == sql.ErrNoRows {
		err = nil
	}
	return days, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBsWeekDays(bsID int) (weekDays []aggregate.WeekDay, err error) {
	command := `SELECT wed.id,
       				   wed.nome as name,
						bwd.start_time,
						bwd.end_time
					FROM weekday wed
					INNER JOIN bs_week_day bwd ON wed.id = bwd.week_day_id
					WHERE bwd.bs_id=$1
					ORDER BY wed.sequence`
	err = aubr.Conn.Select(&weekDays, command, bsID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return weekDays, err
}

func (aubr AppUserBSRepository) GetBsIDServiceIDWeekDays(bsID, serviceId int) (weekDays []entity.WeekDayBS, err error) {
	command := `SELECT wed.id,
						bwd.start_time,
						bwd.end_time,
						bwd.week_day_id
					FROM weekday wed
					INNER JOIN bs_bm_week_day bwd ON wed.id = bwd.week_day_id
					INNER JOIN bs_bm bbm on bwd.bs_bm_id = bbm.id
					INNER JOIN bs_bm_service bbs on bbm.id = bbs.bs_bm_id
					WHERE bbm.bs_id=$1
					AND  bbs.bs_service_id=$2
					GROUP BY wed.id,
						bwd.start_time,
						bwd.end_time,
						bwd.week_day_id,
						wed.sequence
					ORDER BY wed.sequence`
	err = aubr.Conn.Select(&weekDays, command, bsID, serviceId)
	if err == sql.ErrNoRows {
		err = nil
	}
	return weekDays, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSsWithLocalization(criteria criteria.Criteria, lat float64, lon float64) (bss []aggregate.BSAggregate, err error) {
	command := criteria.Query(`SELECT 
                       ST_DISTANCE(bs.geog, ST_SetSRID(ST_MakePoint($1,$2), 4326))   "bs.distance",
                       bs.id,
					   com.cnpj,
					   com.company_name "company.company_name",
					   com.phone  		"company.phone",
					   com.fantasy_name "company.fantasy_name",
					   com.email		"company.email",
					   com.state		"company.state",
					   com.city			"company.city",
					   com.district		"company.district",
					   com.street		"company.street",
					   com.number		"company.number",
					   com.address_complement "company.address_complement",
					   com.zipcode			  "company.zipcode",
					   com.desc				  "company.desc" ,
					   bs.lat				  "bs.lat",
					   bs.lon				  "bs.lon"
                FROM bs
                INNER JOIN company com ON bs.company_id = com.id
				ORDER BY "bs.distance"`)
	err = aubr.Conn.Select(&bss, command, lon, lat)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bss, aubr.Error(err)
}

func (aubr AppUserBSRepository) GetBSsByCategoryWithLocalization(categoryID int, criteria criteria.Criteria, lat float64, lon float64) (bss []aggregate.BSAggregate, err error) {
	command := criteria.Query(`SELECT 
                       ST_DISTANCE(bs.geog, ST_SetSRID(ST_MakePoint($2,$3), 4326)) "bs.distance",
                       bs.id,
					   com.cnpj,
					   com.company_name "company.company_name",
					   com.phone  		"company.phone",
					   com.fantasy_name "company.fantasy_name",
					   com.email		"company.email",
					   com.state		"company.state",
					   com.city			"company.city",
					   com.district		"company.district",
					   com.street		"company.street",
					   com.number		"company.number",
					   com.address_complement "company.address_complement",
					   com.zipcode			  "company.zipcode",
					   com.desc				  "company.desc" ,
					   bs.lat				  "bs.lat",
					   bs.lon				  "bs.lon"
					FROM bs
					INNER JOIN company com ON bs.company_id = com.id
					INNER JOIN bs_bs_category bbc ON bs.id = bbc.bs_id
					WHERE bbc.bs_category_id=$1
					ORDER BY "bs.distance"`)

	err = aubr.Conn.Select(&bss, command, categoryID, lon, lat)
	if err == sql.ErrNoRows {
		err = nil
	}
	return bss, aubr.Error(err)
}
