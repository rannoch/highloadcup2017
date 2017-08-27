package storage

import (
	"github.com/rannoch/highloadcup2017/tcp_server/models"
	"time"
)

var UserDb []*models.User
var LocationDb []*models.Location
var VisitDb []*models.Visit

var UserCount int64
var LocationCount int64
var VisitCount int64

var UserBytesDb [][]byte
var LocationBytesDb [][]byte
var VisitBytesDb [][]byte

var GenerateTime time.Time

// todo later
func InitMemoryMap() {
	UserDb = make([]*models.User, 2000000)
	LocationDb = make([]*models.Location, 2000000)
	VisitDb = make([]*models.Visit, 20000000)

	UserBytesDb = make([][]byte, 2000000)
	LocationBytesDb = make([][]byte, 2000000)
	VisitBytesDb = make([][]byte, 20000000)
}
