package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
	"github.com/rannoch/highloadcup2017/models"
	"strings"
)

var db *sql.DB

func Init() {
	dataSourceName := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=true", "root", "1234", "hlcup2017")
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

func InsertEntity(entity models.Entity) (err error) {
	var tableName string
	var query string
	var questionMarks string

	tableName = entity.TableName()

	for i := 0; i < len(entity.GetValues()); i ++ {
		if i == len(entity.GetValues()) - 1 {
			questionMarks += "?"
			continue
		}

		questionMarks += "?,"
	}

	query = fmt.Sprintf("insert into %s values (%s)", tableName, questionMarks)

	stmtIns, err := db.Prepare(query)

	if err != nil {
		return
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(entity.GetValues()...)


	log.Printf("%v", entity.GetValues())

	return
}

func InsertEntityMultiple(entities []models.Entity) (err error){
	var tableName string
	var query string
	var questionMarks string
	var values = []interface{}{}

	tableName = entities[0].TableName()

	for i := 0; i < len(entities[0].GetValues()); i ++ {
		if i == len(entities[0].GetValues()) - 1 {
			questionMarks += "?"
			continue
		}

		questionMarks += "?,"
	}

	query = fmt.Sprintf("insert into %s values", tableName)

	for _, entity := range entities {
		query += "(" + questionMarks + "),"
		values = append(values, entity.GetValues()...)
	}

	query = strings.TrimSuffix(query, ",")

	stmtIns, err := db.Prepare(query)

	if err != nil {
		return
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(values...)

	return
}

func UpdateEntity(entity models.Entity) (err error) {
	return
}

func ClearAll() (err error){
	return
}