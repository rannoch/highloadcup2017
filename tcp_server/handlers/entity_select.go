package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/tcp_server/storage"
	//"os"
	//"log"
	//"runtime/pprof"
	"bytes"
	"github.com/rannoch/highloadcup2017/tcp_server/server"
)
var UsersNewBytes = []byte(`users/new`)
var LocationsNewBytes = []byte("locations/new")
var VisitsNewBytes = []byte("visits/new")

var UsersBytes = []byte("users")
var LocationsBytes = []byte("locations")
var VisitsBytes = []byte("visits")

var AvgBytes = []byte("avg")
var EmptyJson = []byte(`{}`)

func EntitySelectHandler(ctx *server.HlcupCtx, id int64, entityValue []byte) {
	switch {
	case bytes.Equal(entityValue, UsersBytes):
		if id > storage.UserCount {
			ctx.Error(fasthttp.StatusNotFound)
			return
		}

		ctx.WriteString(storage.UserDb[id].GetString())
		//bufPool.Put(buffer)
	case bytes.Equal(entityValue, LocationsBytes):
		if id > storage.LocationCount {
			ctx.Error(fasthttp.StatusNotFound)
			return
		}

		ctx.WriteString(storage.LocationDb[id].GetString())
		//ctx.Write(storage.LocationBytesDb[id])
	case bytes.Equal(entityValue, VisitsBytes):
		if id > storage.VisitCount {
			ctx.Error(fasthttp.StatusNotFound)
			return
		}

		ctx.WriteString(storage.VisitDb[id].GetString())
		//ctx.Write(storage.VisitBytesDb[id])
	}
}
