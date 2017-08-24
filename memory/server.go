package main

import (
	"github.com/valyala/fasthttp"
	//"github.com/rannoch/highloadcup2017/memory/models"
	"github.com/rannoch/highloadcup2017/memory/handlers"
	"os"
	"log"
	"github.com/rannoch/highloadcup2017/memory/storage"
	//"runtime/pprof"
	//"time"
	"bytes"
	//"fmt"
	"strconv"
	"time"
	"flag"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "/home/baska/projects/go/src/github.com/rannoch/highloadcup2017/memory/memory.prof", "write cpu profile to file")

func main() {
	if len(os.Args) < 3 {
		log.Fatal("not enough args")
	}

	//storage.Init()

	storage.InitMemoryMap()
	LoadData(os.Args[2])

	//flag.Parse()
	if false && *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)

		go func() {
			select {
			case <-time.After(120 * time.Second):
				pprof.StopCPUProfile()
				f.Close()
			}
		}()
		//defer f.Close()
		//defer pprof.StopCPUProfile()
	}

	m := func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		path = path[1:]

		if ctx.IsPost() {
			//POST /<entity>/new на создание

			if bytes.Equal(path, handlers.UsersNewBytes) || bytes.Equal(path, handlers.LocationsNewBytes) || bytes.Equal(path, handlers.VisitsNewBytes) {
				switch {
				case bytes.Equal(path, handlers.UsersNewBytes):
					handlers.EntitityNewHandler(ctx, "users")
				case bytes.Equal(path, handlers.LocationsNewBytes):
					handlers.EntitityNewHandler(ctx, "locations")
				case bytes.Equal(path, handlers.VisitsNewBytes):
					handlers.EntitityNewHandler(ctx, "visits")
				}
				return
			}

			params := bytes.Split(path, []byte("/"))
			if len(params) != 2 {
				ctx.Error("", fasthttp.StatusBadRequest)
				return
			}

			entity := params[0]
			idValue := params[1]

			if !(bytes.Equal(entity, handlers.UsersBytes) || bytes.Equal(entity, handlers.LocationsBytes) || bytes.Equal(entity, handlers.VisitsBytes)) {
				ctx.Error("", fasthttp.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(string(idValue[:]))

			if err != nil || id < 0{
				ctx.Error("", fasthttp.StatusNotFound)
				return
			}

			//POST /<entity>/<id> на обновление
			handlers.EntityUpdateHandler(ctx, int64(id), entity)
			return
		}

		if ctx.IsGet() {
			params := bytes.Split(path, []byte("/"))
			if len(params) < 2 || len(params) > 3 {
				ctx.Error("", fasthttp.StatusBadRequest)
				return
			}

			entity := params[0]
			idValue := params[1]

			if !(bytes.Equal(entity, handlers.UsersBytes) || bytes.Equal(entity, handlers.LocationsBytes) || bytes.Equal(entity, handlers.VisitsBytes)) {
				ctx.Error("", fasthttp.StatusNotFound)
				return
			}

			id, err := strconv.Atoi(string(idValue[:]))
			if err != nil || id < 0{
				ctx.Error("", fasthttp.StatusNotFound)
				return
			}

			if len(params) == 3 {
				if bytes.Equal(entity, handlers.UsersBytes) && bytes.Equal(params[2], handlers.VisitsBytes) {
					handlers.UsersVisitsHandler(ctx, int64(id))
					return
				}
				if bytes.Equal(entity, handlers.LocationsBytes) && bytes.Equal(params[2], handlers.AvgBytes) {
					handlers.LocationsAvgHandler(ctx, int64(id))
					return
				}

				ctx.Error("", fasthttp.StatusBadRequest)
				return
			}

			handlers.EntitySelectHandler(ctx, int64(id), entity)
			return
		}

		ctx.Error("not found", fasthttp.StatusNotFound)
	}

	err := fasthttp.ListenAndServe(":"+os.Args[1], m)

	if err != nil {
		log.Fatal(err.Error())
	}
}
