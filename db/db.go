package db

import (
	"github.com/busgo/pink/conf"
	"github.com/busgo/pink/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DS *sqlx.DB
}

func NewDB(app *conf.AppConf) *DB {

	db, err := sqlx.Connect("mysql", app.Mysql.DSN)
	if err != nil {
		log.Errorf("")
		panic(err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return &DB{DS: db}
}
