package storage

import "github.com/rannoch/highloadcup2017/memory/models"

var UserDb [300000]*models.User
var LocationDb [300000]*models.Location
var VisitDb [2000000]*models.Visit

var UserCount int32
var LocationCount int32
var VisitCount int32

var UserBytesDb [300000][]byte
var LocationBytesDb [300000][]byte
var VisitBytesDb [2000000][]byte

// todo later
func InitMemoryMap() {
	UserDb = [300000]*models.User{}
	LocationDb = [300000]*models.Location{}
	VisitDb = [2000000]*models.Visit{}
}
