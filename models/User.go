package models

import "time"

type User struct {
	Id int32
	Email string
	first_name string
	last_name string
	gender string
	birth_date time.Time
}