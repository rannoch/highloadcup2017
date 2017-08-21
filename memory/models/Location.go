package models

import (
	"github.com/rannoch/highloadcup2017/util"
	"reflect"
	"strings"
	"fmt"
)

type Location struct {
	Id       int32 `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance int32 `json:"distance"`
	Visits   []*Visit `json:"-"`
}

func (location *Location) HasForeignRelations() bool {
	return false
}

func (location *Location) TableName() string {
	return "location"
}

func (location *Location) GetId() int32 {
	return location.Id
}

func (location *Location) GetFields(alias string) []string {
	return []string{"id", "place", "country", "city", "distance"}
}

func (location *Location) GetValues() []interface{} {
	return []interface{}{location.Id, location.Place, location.Country, location.City, location.Distance}
}

func (location *Location) GetFieldPointers(with []string) []interface{} {
	return []interface{}{&location.Id, &location.Place, &location.Country, &location.City, &location.Distance}
}

func (location *Location) ValidateParams(params map[string]interface{}, scenario string) (result bool) {
	if scenario == "insert" && len(params) != len(location.GetFields("")) {
		return false
	}

	for param, value := range params {
		if value == nil {
			return false
		}

		if scenario == "update" && param == "id" {
			return false
		}

		if !util.StringInSlice(param, location.GetFields("")) {
			return false
		}
	}

	return true
}

func (location *Location) SetParams(params map[string]interface{}) {
	locationValue := reflect.ValueOf(location).Elem()

	for param, value := range params {
		field := locationValue.FieldByName(strings.Title(param))

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

func (location *Location) GetBytes() []byte {
	return []byte(fmt.Sprintf("{\"id\":%d,\"place\":\"%s\",\"country\":\"%s\",\"city\":\"%s\",\"distance\":%d}", location.Id, location.Place, location.Country, location.City, location.Distance))
}
