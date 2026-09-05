package contracts

import (
	"bitbucket.org/lyndus/backend/infra/db"
	"github.com/jmoiron/sqlx"
)

type RepositoryInterface interface {
	GetConnector() db.Connector
	GetTable() string
	GetConnection() *sqlx.DB
}