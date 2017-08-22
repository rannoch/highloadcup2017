package models

import (
	"github.com/rannoch/highloadcup2017/mysql_server/util"
	"fmt"
	"encoding/json"
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
	id, ok := params["id"]; if ok {
		switch id.(type) {
		case int32:
			location.Id = id.(int32)
		case float32:
			location.Id = int32(id.(float32))
		case float64:
			location.Id = int32(id.(float64))
		case json.Number:
			t, _ := id.(json.Number).Int64()
			location.Id = int32(t)
		}
	}
	place, ok := params["place"]; if ok {
		location.Place = place.(string)
	}
	country, ok := params["country"]; if ok {
		location.Country = country.(string)
	}
	city, ok := params["city"]; if ok {
		location.City = city.(string)
	}
	distance, ok := params["distance"]; if ok {
		switch distance.(type) {
		case int32:
			location.Distance = distance.(int32)
		case float32:
			location.Distance = int32(distance.(float32))
		case float64:
			location.Distance = int32(distance.(float64))
		case json.Number:
			t, _ := distance.(json.Number).Int64()
			location.Distance = int32(t)
		}
	}
}

func (location *Location) GetBytes() []byte {
	return []byte(fmt.Sprintf("{\"id\":%d,\"place\":\"%s\",\"country\":\"%s\",\"city\":\"%s\",\"distance\":%d}", location.Id, location.Place, location.Country, location.City, location.Distance))
}

func (location *Location) GetString() string {
	return fmt.Sprintf("{\"id\":%d,\"place\":\"%s\",\"country\":\"%s\",\"city\":\"%s\",\"distance\":%d}", location.Id, location.Place, location.Country, location.City, location.Distance)
}