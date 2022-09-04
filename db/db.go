package db

import (
	"context"
	"database/sql"
	"fmt"
	"godok/config"
	"log"

	_ "github.com/lib/pq"
)

var ApiDb *sql.DB

type Queries struct {
	db DBTX
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func Connect() *Queries {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	)
	var err error
	ApiDb, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	q := Queries{db: ApiDb}
	return &q
}

func Close() {
	ApiDb.Close()
}
