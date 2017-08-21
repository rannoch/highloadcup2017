package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"flag"
	//"os"
	//"log"
	//"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "/home/baska/projects/go/src/github.com/rannoch/highloadcup2017/memory/memory.prof", "write cpu profile to file")

func EntitySelectHandler(ctx *fasthttp.RequestCtx) {
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

	var entityValue string
	var id int32

	id, _ = ctx.UserValue("id").(int32)

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok || !(entityValue == "users" || entityValue == "locations" || entityValue == "visits"){
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	switch entityValue {
	case "users":
		if id > storage.UserCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := storage.UserDb[id]

		if entity == nil {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		ctx.SetBody(entity.GetBytes())
	case "locations":
		if id > storage.LocationCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := storage.LocationDb[id]

		if entity == nil {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		ctx.SetBody(entity.GetBytes())
	case "visits":
		if id > storage.VisitCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := storage.VisitDb[id]

		if entity == nil {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		ctx.SetBody(entity.GetBytes())
	}
}
