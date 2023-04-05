package db

import (
	"database/sql"

	"github.com/idprm/go-pass-tsel/src/config"
	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Secret) *sql.DB {
	dsn := cfg.Db.Source

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db
}
