package storage

import "github.com/rannoch/highloadcup2017/memory/models"

var UserDb []*models.User
var LocationDb []*models.Location
var VisitDb []*models.Visit

var UserCount int64
var LocationCount int64
var VisitCount int64

var UserBytesDb [][]byte
var LocationBytesDb [][]byte
var VisitBytesDb [][]byte

// todo later
func InitMemoryMap() {
	UserDb = make([]*models.User, 200000)
	LocationDb = make([]*models.Location, 200000)
	VisitDb = make([]*models.Visit, 2000000)

	UserBytesDb = make([][]byte, 200000)
	LocationBytesDb = make([][]byte, 200000)
	VisitBytesDb = make([][]byte, 2000000)
}
