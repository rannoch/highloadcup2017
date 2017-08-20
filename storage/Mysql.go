package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"reflect"
)

type mysqlDb struct {
	*sql.DB
}

func Init() {
	dataSourceName := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=true", "root", "1234", "hlcup2017")
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

	mysqlDb.SetMaxOpenConns(200)

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

	return
}

func (db mysqlDb) InsertEntityMultiple(entities []Entity) (err error) {
	if len(entities) < 1 {
		return
	}

	var tableName string
	var query string
	var valuesString string

	tableName = entities[0].TableName()

	for _, entity := range entities {
		s := ""

		for _, entityFieldValue := range entity.GetValues() {
			switch entityFieldValue.(type) {
			case string:
				s += fmt.Sprintf("'%s',", entityFieldValue)
			default:
				s += fmt.Sprintf("%v,", entityFieldValue)
			}
		}

		s = strings.TrimSuffix(s, ",")

		valuesString += "(" + s + "),"
	}

	valuesString = strings.TrimSuffix(valuesString, ",")

	query = fmt.Sprintf("insert into %s values %s", tableName, valuesString)

	_, err = db.Exec(query)

	return
}

func (db mysqlDb) UpdateEntity(entity Entity, params map[string]interface{}, conditions []Condition) (rowsAffected int64, err error) {
	if len(params) < 1 {
		return
	}

	var query string
	var conditionString string
	var setString string
	var values []interface{}

	if len(conditions) > 0 {
		conditionString += "where "
	}
	for i := 0; i < len(conditions); i++ {
		condition := conditions[i]

		if i > 0 {
			conditionString += condition.JoinCondition + " "
		}

		conditionString += fmt.Sprintf("%s %s %s", condition.Param, condition.Operator, condition.Value)
	}

	var i = 0
	for param, value := range params {
		if i > 0 {
			setString += ","
		}

		setString += fmt.Sprintf("%s = ?", param)
		values = append(values, value)

		i++
	}

	query = fmt.Sprintf("update %s set %s %s", entity.TableName(), setString, conditionString)

	//log.Println(query)

	stmtIns, err := db.Prepare(query)

	if err != nil {
		return
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(values...)

	return result.RowsAffected()
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
			conditionString += condition.JoinCondition + " "
		}

		conditionString += fmt.Sprintf("%s %s %s", condition.Param, condition.Operator, condition.Value)
	}

	query = fmt.Sprintf("select * from %s %s limit 1", entity.TableName(), conditionString)

	err = db.QueryRow(query).Scan(entity.GetFieldPointers([]string{})...)

	return
}

func (db mysqlDb) SelectEntityMultiple(out interface{}, fields []string, tableJoins []Join, conditions []Condition, sort Sort) (err error) {
	var query string
	var conditionString string
	var joinString string
	var fieldsString string
	var with []string
	var sortString string

	entityType := reflect.TypeOf(out).Elem().Elem()
	outValue := reflect.Indirect(reflect.ValueOf(out))

	//fmt.Println(entityType)
	reflectEntity := (reflect.New(entityType).Interface()).(Entity)

	tableName := reflectEntity.TableName()

	// fields
	if len(fields) < 1 {
		fieldsString = "*"
	} else {
		fieldsString = strings.Join(fields, ",")
	}

	// join
	for i := 0; i < len(tableJoins); i++ {
		tableJoin := tableJoins[i]
		with = append(with, tableJoin.Name)

		joinString += tableJoin.Type + " join " + tableJoin.Name + " on " + fmt.Sprintf("%s %s %s", tableJoin.Condition.Param, tableJoin.Condition.Operator, tableJoin.Condition.Value)
	}

	// where
	if len(conditions) > 0 {
		conditionString += "where "
	}
	for i := 0; i < len(conditions); i++ {
		condition := conditions[i]

		if i > 0 {
			conditionString += condition.JoinCondition + " "
		}

		conditionString += fmt.Sprintf("%s %s %s ", condition.Param, condition.Operator, condition.Value)
	}

	// sort
	if len(sort.Fields) > 0 {
		sortString += " order by "

		for _, sortField := range sort.Fields {
			sortString += " " + sortField
		}

		sortString += " " + sort.Direction + " "
	}

	query = fmt.Sprintf("select %s from %s %s %s %s", fieldsString, tableName, joinString, conditionString, sortString)

	//fmt.Println(query)

	rows, err := db.Query(query)

	if err != nil {
		return
	}

	for rows.Next() {
		var entity Entity = (reflect.New(entityType).Interface()).(Entity)

		err = rows.Scan(entity.GetFieldPointers(with)...)

		if err != nil {
			return
		}

		outValue.Set(reflect.Append(outValue, reflect.ValueOf(entity).Elem()))
	}

	return
}

func (db mysqlDb) GetAverage(entity Entity, avgColumn string, tableJoins []Join, conditions []Condition) (average float32, err error){
	var query string
	var conditionString string
	var joinString string

	// join
	for i := 0; i < len(tableJoins); i++ {
		tableJoin := tableJoins[i]

		joinString += tableJoin.Type + " join " + tableJoin.Name + " on " + fmt.Sprintf("%s %s %s", tableJoin.Condition.Param, tableJoin.Condition.Operator, tableJoin.Condition.Value)
	}

	// where
	if len(conditions) > 0 {
		conditionString += "where "
	}
	for i := 0; i < len(conditions); i++ {
		condition := conditions[i]

		if i > 0 {
			conditionString += condition.JoinCondition + " "
		}

		conditionString += fmt.Sprintf("%s %s %s ", condition.Param, condition.Operator, condition.Value)
	}

	query = fmt.Sprintf("select round(avg(%s), 5) from %s %s %s", avgColumn, entity.TableName(), joinString, conditionString)

	//fmt.Println(query)

	err = db.QueryRow(query).Scan(&average)

	return
}