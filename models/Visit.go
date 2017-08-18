package models

type Visit struct {
	Id int32 `json:"id"`
	Location int32 `json:"location"`
	User int32 `json:"user"`
	Visited_at BaskaTime `json:"visited_at"`
	Mark int32 `json:"mark"`
}

func (visit *Visit) TableName() string {
	return "visit"
}

func (visit *Visit) GetId() int32 {
	return visit.Id
}

func (visit *Visit) GetFields() []string{
	return []string{"id", "location", "user", "visited_at", "mark"}
}

func (visit *Visit) GetValues() []interface{} {
	return []interface{}{visit.Id, visit.Location, visit.User, visit.Visited_at, visit.Mark}
}

func (visit *Visit) GetFieldPointers() []interface{} {
	return []interface{}{&visit.Id, &visit.Location, &visit.User, &visit.Visited_at, &visit.Mark}
}