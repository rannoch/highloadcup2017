package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"reflect"
)

var Db IStorage

type mysqlDb struct {
	*sql.DB
}

func Init() {
	dataSourceName := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=true", "root", "1", "hlcup2017")
	var err error
	var mysqlDb mysqlDb

	mysqlDb.DB, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	err = mysqlDb.Ping()

	if err != nil {
		log.Fatal(err)
	}

	Db = mysqlDb
}

func (db mysqlDb) InsertEntity(entity Entity) (err error) {
	var tableName string
	var query string
	var questionMarks string

	tableName = entity.TableName()

	for i := 0; i < len(entity.GetValues()); i ++ {
		if i == len(entity.GetValues())-1 {
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

func (db mysqlDb) InsertEntityMultiple(entities []Entity) (err error) {
	var tableName string
	var query string
	var questionMarks string
	var values = []interface{}{}

	tableName = entities[0].TableName()

	for i := 0; i < len(entities[0].GetValues()); i ++ {
		if i == len(entities[0].GetValues())-1 {
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

func UpdateEntity(entity Entity) (err error) {
	return
}

func ClearAll() (err error) {
	return
}

func (db mysqlDb) SelectEntity(entity Entity, conditions []Condition) (err error) {
	var query string
	var conditionString string

	if len(conditions) > 0 {
		conditionString += "where "
	}
	for i := 0; i < len(conditions); i++ {
		condition := conditions[i]

		if i > 0 {
			conditionString += condition.Join + " "
		}

		conditionString += fmt.Sprintf("%s %s %s", condition.Param, condition.Operator, condition.Value)
	}

	query = fmt.Sprintf("select * from %s %s limit 1", entity.TableName(), conditionString)

	err = db.QueryRow(query).Scan(entity.GetFieldPointers()...)

	return
}

func (db mysqlDb) SelectEntityMultiple(out interface{}, conditions []Condition) (err error) {
	var entities []Entity

	var query string
	var conditionString string

	entityType := reflect.TypeOf(out).Elem().Elem()
	outValue := reflect.Indirect(reflect.ValueOf(out))

	fmt.Println(entityType)
	reflectEntity := (reflect.New(entityType).Interface()).(Entity)

	tableName := reflectEntity.TableName()

	if len(conditions) > 0 {
		conditionString += "where "
	}
	for i := 0; i < len(conditions); i++ {
		condition := conditions[i]

		if i > 0 {
			conditionString += condition.Join + " "
		}

		conditionString += fmt.Sprintf("%s %s %s", condition.Param, condition.Operator, condition.Value)
	}

	query = fmt.Sprintf("select * from %s %s limit 1", tableName, conditionString)

	rows, err := db.Query(query)

	if err != nil {
		return
	}

	for rows.Next() {
		var entity Entity = reflectEntity

		err = rows.Scan(entity.GetFieldPointers()...)

		fmt.Println(entity)
		if err != nil {
			return
		}

		entities = append(entities, entity)
	}


	outValue.Set(reflect.ValueOf(entities))
	return
}
