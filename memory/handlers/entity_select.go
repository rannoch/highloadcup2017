package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	//"os"
	//"log"
	//"runtime/pprof"
	"sync"
	"bytes"
)
var UsersNewBytes = []byte(`users/new`)
var LocationsNewBytes = []byte("locations/new")
var VisitsNewBytes = []byte("visits/new")

var UsersBytes = []byte("users")
var LocationsBytes = []byte("locations")
var VisitsBytes = []byte("visits")

var AvgBytes = []byte("avg")


var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

func EntitySelectHandler(ctx *fasthttp.RequestCtx, id int32, entityValue []byte) {
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

		entity := storage.UserDb[id]

		if entity == nil {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		buffer := bufPool.Get().(*bytes.Buffer)
		buffer.Reset()
		buffer.WriteString(entity.GetString())

		ctx.Write(buffer.Bytes())
		bufPool.Put(buffer)
	case bytes.Equal(entityValue, LocationsBytes):
		if id > storage.LocationCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := storage.LocationDb[id]

		if entity == nil {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		buffer := bufPool.Get().(*bytes.Buffer)
		buffer.Reset()
		buffer.WriteString(entity.GetString())

		ctx.Write(buffer.Bytes())
		bufPool.Put(buffer)
	case bytes.Equal(entityValue, VisitsBytes):
		if id > storage.VisitCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := storage.VisitDb[id]

		if entity == nil {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		buffer := bufPool.Get().(*bytes.Buffer)
		buffer.Reset()
		buffer.WriteString(entity.GetString())

		ctx.Write(buffer.Bytes())
		bufPool.Put(buffer)
	}
}
