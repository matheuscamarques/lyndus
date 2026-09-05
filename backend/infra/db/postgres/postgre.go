package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

//InitDB database init
func InitDB(host, user, pass, database string, port int) (db *sqlx.DB, err error) {
	dbInfo := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", user, pass, host, port, database)
	db, err = sqlx.Open("postgres", dbInfo)
	if err != nil {
		return
	}
	//db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return
}
