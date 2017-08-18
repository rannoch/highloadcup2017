package models

type Location struct {
	Id int32 `json:"id"`
	Place string `json:"place"`
	Country string `json:"country"`
	City string `json:"city"`
	Distance int32 `json:"distance"`
}

func (location *Location) TableName() string {
	return "location"
}

func (location *Location) GetId() int32 {
	return location.Id
}

func (location *Location) GetFields() []string{
	return []string{"id", "place", "country", "city", "distance"}
}

func (location *Location) GetValues() []interface{} {
	return []interface{}{location.Id, location.Place, location.Country, location.City, location.Distance}
}

func (location *Location) GetFieldPointers() []interface{} {
	return []interface{}{&location.Id, &location.Place, &location.Country, &location.City, &location.Distance}
}