package db

import (
	"github.com/busgo/pink/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
}

func NewDB(app *conf.AppConf) *DB {

	db, err := sqlx.Connect("mysql", "")
	if err != nil {
		panic(err)
	}
	return &DB{db: db}
}
