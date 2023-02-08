package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	Client *sql.DB

	username = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASS")
	host     = os.Getenv("DB_HOST")
	schema   = os.Getenv("DB_SCHEMA")
)

func init() {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		username,
		password,
		host,
		schema,
	)

	var err error
	Client, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	log.Println("Connected to database successfully")
}
