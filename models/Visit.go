package models

type Visit struct {
	Id int32
	Location int32
	User int32
	Visited_at BaskaTime
	Mark int32
}

func (visit Visit) TableName() string {
	return "visit"
}

func (visit Visit) GetId() int32 {
	return visit.Id
}

func (Visit Visit) GetFields() []string{
	return []string{"id", "location", "user", "visited_at", "mark"}
}

func (Visit Visit) GetValues() []interface{} {
	return []interface{}{Visit.Id, Visit.Location, Visit.User, Visit.Visited_at, Visit.Mark}
}