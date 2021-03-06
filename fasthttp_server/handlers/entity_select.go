package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	//"os"
	//"log"
	//"runtime/pprof"
	"bytes"
)
var UsersNewBytes = []byte(`users/new`)
var LocationsNewBytes = []byte("locations/new")
var VisitsNewBytes = []byte("visits/new")

var UsersBytes = []byte("users")
var LocationsBytes = []byte("locations")
var VisitsBytes = []byte("visits")

var AvgBytes = []byte("avg")
var EmptyJson = []byte(`{}`)

func EntitySelectHandler(ctx *fasthttp.RequestCtx, id int64, entityValue []byte) {
	/*flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}*/

	ctx.SetContentType("application/json;charset=utf-8")

	switch {
	case bytes.Equal(entityValue, UsersBytes):
		if id > storage.UserCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		//buffer := bufPool.Get().(*bytes.Buffer)
		//buffer.Reset()
		//buffer.Write(storage.UserBytesDb[id])

		//ctx.Write(storage.UserBytesDb[id])
		ctx.SetBody(storage.UserBytesDb[id])
		//bufPool.Put(buffer)
	case bytes.Equal(entityValue, LocationsBytes):
		if id > storage.LocationCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		//buffer := bufPool.Get().(*bytes.Buffer)
		//buffer.Reset()
		//buffer.Write(storage.LocationBytesDb[id])

		//ctx.Write(storage.LocationBytesDb[id])
		ctx.SetBody(storage.LocationBytesDb[id])
		//bufPool.Put(buffer)
	case bytes.Equal(entityValue, VisitsBytes):
		if id > storage.VisitCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		//buffer := bufPool.Get().(*bytes.Buffer)
		//buffer.Reset()
		//buffer.Write(storage.VisitBytesDb[id])

		//ctx.Write(storage.VisitBytesDb[id])
		ctx.SetBody(storage.VisitBytesDb[id])
		//bufPool.Put(buffer)
	}
}
