package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func Init() {
	dataSourceName := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=true", "root", "1", "hlcup2017" )
	var err error

	db, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}
}