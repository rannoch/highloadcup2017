package models

import (
	"github.com/rannoch/highloadcup2017/mysql_server/util"
	"reflect"
	"strings"
)

type Visit struct {
	Id            int32 `json:"id"`
	Location      int32 `json:"location"`
	User          int32 `json:"user"`
	Visited_at    int32 `json:"visited_at"`
	Mark          int32 `json:"mark"`
	LocationChild *Location `json:"-" relation:"location"`
	UserChild     *User `json:"-" relation:"user"`
}

func (visit *Visit) HasForeignRelations() bool {
	return true
}

func (visit *Visit) TableName() string {
	return "visit"
}

func (visit *Visit) GetId() int32 {
	return visit.Id
}

func (visit *Visit) GetFields(alias string) []string {
	if len(alias) > 0 {
		return []string{alias + ".id", alias + ".location", alias + ".user", alias + ".visited_at", alias + ".mark"}
	}

	return []string{"id", "location", "user", "visited_at", "mark"}
}

func (visit *Visit) ValidateParams(params map[string]interface{}, scenario string) (result bool) {
	if scenario == "insert" && len(params) != len(visit.GetFields("")) {
		return false
	}

	for param, value := range params {
		if value == nil {
			return false
		}

		if scenario == "update" && param == "id" {
			return false
		}

		if !util.StringInSlice(param, visit.GetFields("")) {
			return false
		}
	}

	return true
}

func (visit *Visit) GetValues() []interface{} {
	return []interface{}{visit.Id, visit.Location, visit.User, visit.Visited_at, visit.Mark}
}

func (visit *Visit) GetFieldPointers(with []string) []interface{} {
	fieldPointers := []interface{}{&visit.Id, &visit.Location, &visit.User, &visit.Visited_at, &visit.Mark}

	for _, v := range with {
		if v == "location" {
			if visit.LocationChild == nil {
				visit.LocationChild = &Location{}
			}

			fieldPointers = append(fieldPointers, visit.LocationChild.GetFieldPointers([]string{})...)
		}
	}

	return fieldPointers
}

func (visit *Visit) SetParams(params map[string]interface{}) {
	visitValue := reflect.ValueOf(visit).Elem()

	for param, value := range params {
		field := visitValue.FieldByName(strings.Title(param))

		switch field.Interface().(type) {
		case int32:
			switch value.(type) {
			case int32:
				field.Set(reflect.ValueOf(value.(int32)))
			case float32:
				field.Set(reflect.ValueOf(int32(value.(float32))))
			case float64:
				field.Set(reflect.ValueOf(int32(value.(float64))))
			}
		case string:
			field.SetString(value.(string))
		}
	}
}
