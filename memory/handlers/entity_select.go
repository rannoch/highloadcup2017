package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"strconv"
	"flag"
	//"os"
	//"log"
	//"runtime/pprof"
	"github.com/rannoch/highloadcup2017/memory/models"
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
	var id int
	var entity interface{}

	id, err := strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok || !(entityValue == "users" || entityValue == "locations" || entityValue == "visits"){
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	entity = storage.Db[entityValue[:len(entityValue) - 1]][int32(id)]

	if entity == nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	switch entity.(type) {
	case *models.Location:
		ctx.SetBody(entity.(*models.Location).GetBytes())
	case *models.User:
		ctx.SetBody(entity.(*models.User).GetBytes())
	case *models.Visit:
		ctx.SetBody(entity.(*models.Visit).GetBytes())
	}
}
