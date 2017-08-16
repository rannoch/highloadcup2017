package models

import "time"

type Visit struct {
	Id int32
	Location int32
	User int32
	Visited_at time.Time
	Mark int32
}
