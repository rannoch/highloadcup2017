package models

import (
	"github.com/rannoch/highloadcup2017/mysql_server/util"
	"fmt"
	"encoding/json"
)

type Location struct {
	Id       int64
	Place    []byte
	Country  []byte
	City     []byte
	Distance int64
	Visits   []*Visit
}

type LocationHelper struct {
	Id       int64
	Place    string
	Country  string
	City     string
	Distance int64
	Visits   []*Visit
}

func (locationHelper LocationHelper) GetLocation() Location {
	return Location{
		Id:locationHelper.Id,
		Place: []byte(locationHelper.Place),
		Country: []byte(locationHelper.Country),
		City: []byte(locationHelper.City),
		Distance:locationHelper.Distance,
		Visits:locationHelper.Visits,
	}
}

func (location *Location) GetFields(alias string) []string {
	return []string{"id", "place", "country", "city", "distance"}
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
			location.Id = int64(id.(int32))
		case int64:
			location.Id = id.(int64)
		case float32:
			location.Id = int64(id.(float32))
		case float64:
			location.Id = int64(id.(float64))
		case json.Number:
			location.Id, _ = id.(json.Number).Int64()
		}
	}
	place, ok := params["place"]; if ok {
		location.Place = []byte(place.(string))
	}
	country, ok := params["country"]; if ok {
		location.Country = []byte(country.(string))
	}
	city, ok := params["city"]; if ok {
		location.City = []byte(city.(string))
	}
	distance, ok := params["distance"]; if ok {
		switch distance.(type) {
		case int32:
			location.Distance = int64(distance.(int32))
		case int64:
			location.Distance = distance.(int64)
		case float32:
			location.Distance = int64(distance.(float32))
		case float64:
			location.Distance = int64(distance.(float64))
		case json.Number:
			location.Distance, _ = distance.(json.Number).Int64()
		}
	}
}

func (location *Location) GetBytes() []byte {
	return []byte(fmt.Sprintf("{\"id\":%d,\"place\":\"%s\",\"country\":\"%s\",\"city\":\"%s\",\"distance\":%d}", location.Id, location.Place, location.Country, location.City, location.Distance))
}

func (location *Location) GetString() string {
	return fmt.Sprintf("{\"id\":%d,\"place\":\"%s\",\"country\":\"%s\",\"city\":\"%s\",\"distance\":%d}", location.Id, location.Place, location.Country, location.City, location.Distance)
}