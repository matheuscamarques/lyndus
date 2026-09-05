package db

import (
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ConnectorInterface interface {
	SetConn(conn *sqlx.DB)
	UsePostgres()
	Error(err error)error
	GetTable() string
	GetConnector() Connector
	GetConnection() *sqlx.DB
}

type Connector struct {
	Conn  *sqlx.DB
	Table string
}

func (c *Connector) SetConn(conn *sqlx.DB) {
	c.Conn = conn
}

func (c *Connector) UsePostgres() {
	c.Conn = postgres.DB
}

func (c Connector) Error(err error) error{
	if err == nil {
		return err
	}
	return errors.New(c.GetTable() +" "+ err.Error())
}
func (c Connector) GetTable() string{
	return c.Table
}

func (c Connector) GetConnector() Connector {
	return c
}

func (c Connector)GetConnection() *sqlx.DB{
	return c.Conn
}


