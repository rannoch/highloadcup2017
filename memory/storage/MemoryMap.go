package storage

import "github.com/rannoch/highloadcup2017/memory/models"

var UserDb map[int32]*models.User
var LocationDb map[int32]*models.Location
var VisitDb map[int32]*models.Visit

// todo later
func InitMemoryMap() {
	UserDb = map[int32]*models.User{}
	LocationDb = map[int32]*models.Location{}
	VisitDb = map[int32]*models.Visit{}
}
