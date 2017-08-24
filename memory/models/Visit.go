package models

import (
	"github.com/rannoch/highloadcup2017/mysql_server/util"
	"fmt"
	"encoding/json"
)

type Visit struct {
	Id             int64
	Location       int64
	User           int64
	Visited_at     int64
	Mark           int64
	Location_model *Location
	User_model     *User
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

func (visit *Visit) SetParams(params map[string]interface{}) {
	id, ok := params["id"]; if ok {
		switch id.(type) {
		case int32:
			visit.Id = int64(id.(int32))
		case float32:
			visit.Id = int64(id.(float32))
		case float64:
			visit.Id = int64(id.(float64))
		}
	}
	location, ok := params["location"]; if ok {
		switch location.(type) {
		case int32:
			visit.Location = int64(location.(int32))
		case int64:
			visit.Location = location.(int64)
		case float32:
			visit.Location = int64(location.(float32))
		case float64:
			visit.Location = int64(location.(float64))
		case json.Number:
			visit.Location, _ = location.(json.Number).Int64()
		}
	}

	user, ok := params["user"]; if ok {
		switch user.(type) {
		case int32:
			visit.User = int64(user.(int32))
		case int64:
			visit.User = user.(int64)
		case float32:
			visit.User = int64(user.(float32))
		case float64:
			visit.User = int64(user.(float64))
		case json.Number:
			visit.User, _ = user.(json.Number).Int64()
		}
	}

	visited_at, ok := params["visited_at"]; if ok {
		switch visited_at.(type) {
		case int32:
			visit.Visited_at = int64(visited_at.(int32))
		case int64:
			visit.Visited_at = visited_at.(int64)
		case float32:
			visit.Visited_at = int64(visited_at.(float32))
		case float64:
			visit.Visited_at = int64(visited_at.(float64))
		case json.Number:
			visit.Visited_at, _ = visited_at.(json.Number).Int64()
		}
	}

	mark, ok := params["mark"]; if ok {
		switch mark.(type) {
		case int32:
			visit.Mark = int64(mark.(int32))
		case int64:
			visit.Mark = mark.(int64)
		case float32:
			visit.Mark = int64(mark.(float32))
		case float64:
			visit.Mark = int64(mark.(float64))
		case json.Number:
			visit.Mark, _ = mark.(json.Number).Int64()
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

func (visit *Visit) GetString() string {
	return fmt.Sprintf("{\"id\":%d,\"location\":%d,\"user\":%d,\"visited_at\":%d,\"mark\":%d}", visit.Id, visit.Location, visit.User, visit.Visited_at, visit.Mark)
}
