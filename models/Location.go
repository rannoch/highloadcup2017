package models

type Location struct {
	Id int32
	Place string
	Country string
	City string
	Distance int32
}

func (location Location) TableName() string {
	return "location"
}


func (location Location) GetId() int32 {
	return location.Id
}

func (location Location) GetFields() []string{
	return []string{"id", "place", "country", "city", "distance"}
}

func (location Location) GetValues() []interface{} {
	return []interface{}{location.Id, location.Place, location.Country, location.City, location.Distance}
}