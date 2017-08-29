package storage

import (
	"github.com/rannoch/highloadcup2017/epoll_server/models"
	"time"
)

var UserDb [1200000]*models.User
var LocationDb [1200000]*models.Location
var VisitDb [15000000]*models.Visit

var UserCount int64
var LocationCount int64
var VisitCount int64

var UserBytesDb [1200000][]byte
var LocationBytesDb [1200000][]byte
var VisitBytesDb [15000000][]byte

var GenerateTime time.Time

// todo later
func InitMemoryMap() {

}
