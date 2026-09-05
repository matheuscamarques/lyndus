package repository

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type AppUserServiceRepository struct {
	contracts.AppUserServiceRepositoryInterface
	conn *sqlx.DB
}

func NewAppUserServiceRepository() *AppUserServiceRepository {
	return &AppUserServiceRepository{
		conn: postgres.DB,
	}
}

func (aus AppUserServiceRepository) GetServiceCategory() (serviceCategories []entity.ServiceCategory, err error) {
	command := `SELECT id, "desc" FROM bs_service_category ORDER BY id `
	err = aus.conn.Select(&serviceCategories, command)
	if err == sql.ErrNoRows {
		err = nil
	}
	return serviceCategories, err
}

func (aus AppUserServiceRepository) GetServicesByCategory(bsID, serviceCategoryID int) (services []entity.Service, err error) {
	command := `SELECT id,
       				   "name",
       				   "desc",
       				   "value"
					FROM bs_service
					WHERE bs_id=$1 AND bs_service_category_id=$2
					ORDER BY name `
	err = aus.conn.Select(&services, command, bsID, serviceCategoryID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return services, err
}
