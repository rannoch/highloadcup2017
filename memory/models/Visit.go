package models

import (
	"github.com/rannoch/highloadcup2017/util"
	"fmt"
	"encoding/json"
)

type Visit struct {
	Id             int32 `json:"id"`
	Location       int32 `json:"location"`
	User           int32 `json:"user"`
	Visited_at     int32 `json:"visited_at"`
	Mark           int32 `json:"mark"`
	Location_model *Location `json:"-" relation:"location"`
	User_model     *User `json:"-" relation:"user"`
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
			if visit.Location_model == nil {
				visit.Location_model = &Location{}
			}

			fieldPointers = append(fieldPointers, visit.Location_model.GetFieldPointers([]string{})...)
		}
	}

	return fieldPointers
}

func (visit *Visit) SetParams(params map[string]interface{}) {
	id, ok := params["id"]; if ok {
		switch id.(type) {
		case int32:
			visit.Id = id.(int32)
		case float32:
			visit.Id = int32(id.(float32))
		case float64:
			visit.Id = int32(id.(float64))
		}
	}
	location, ok := params["location"]; if ok {
		switch location.(type) {
		case int32:
			visit.Location = location.(int32)
		case float32:
			visit.Location = int32(location.(float32))
		case float64:
			visit.Location = int32(location.(float64))
		case json.Number:
			t, _ := location.(json.Number).Int64()
			visit.Location = int32(t)
		}
	}

	user, ok := params["user"]; if ok {
		switch user.(type) {
		case int32:
			visit.User = user.(int32)
		case float32:
			visit.User = int32(user.(float32))
		case float64:
			visit.User = int32(user.(float64))
		case json.Number:
			t, _ := user.(json.Number).Int64()
			visit.User = int32(t)
		}
	}

	visited_at, ok := params["visited_at"]; if ok {
		switch visited_at.(type) {
		case int32:
			visit.Visited_at = visited_at.(int32)
		case float32:
			visit.Visited_at = int32(visited_at.(float32))
		case float64:
			visit.Visited_at = int32(visited_at.(float64))
		case json.Number:
			t, _ := visited_at.(json.Number).Int64()
			visit.Visited_at = int32(t)
		}
	}

	mark, ok := params["mark"]; if ok {
		switch mark.(type) {
		case int32:
			visit.Mark = mark.(int32)
		case float32:
			visit.Mark = int32(mark.(float32))
		case float64:
			visit.Mark = int32(mark.(float64))
		case json.Number:
			t, _ := mark.(json.Number).Int64()
			visit.Mark = int32(t)
		}
	}

}

type VisitByDateAsc []*Visit

func (s VisitByDateAsc) Len() int {
	return len(s)
}
func (s VisitByDateAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s VisitByDateAsc) Less(i, j int) bool {
	return s[i].Visited_at < s[j].Visited_at
}

func (visit *Visit) GetBytes() []byte {
	return []byte(fmt.Sprintf("{\"id\":%d,\"location\":%d,\"user\":%d,\"visited_at\":%d,\"mark\":%d}", visit.Id, visit.Location, visit.User, visit.Visited_at, visit.Mark))
}
