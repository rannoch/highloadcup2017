package main

import (
	"github.com/valyala/fasthttp"
	//"github.com/rannoch/highloadcup2017/memory/models"
	"github.com/rannoch/highloadcup2017/memory/handlers"
	"os"
	"log"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"strings"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("not enough args")
	}

	//storage.Init()

	storage.InitMemoryMap()
	LoadData(os.Args[2])

	m := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		path = path[1:]

		if ctx.IsPost() {
			//POST /<entity>/new на создание

			if path == "users/new" || path == "locations/new" || path == "visits/new" {
				switch path {
				case "users/new":
					ctx.SetUserValue("entity", "users")
				case "locations/new":
					ctx.SetUserValue("entity", "locations")
				case "visits/new":
					ctx.SetUserValue("entity", "visits")
				}

				handlers.EntitityNewHandler(ctx)
				return
			}

			params := strings.Split(path, "/")
			if len(params) != 2 {
				ctx.Error("", fasthttp.StatusBadRequest)
				return
			}

			entity := params[0]
			idValue := params[1]

			if !(entity == "users" || entity == "locations" || entity == "visits") {
				ctx.Error("", fasthttp.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(idValue)
			if err != nil {
				ctx.Error("", fasthttp.StatusNotFound)
				return
			}

			ctx.SetUserValue("id", int32(id))
			ctx.SetUserValue("entity", entity)

			//POST /<entity>/<id> на обновление
			handlers.EntityUpdateHandler(ctx)
			return
		}

		if ctx.IsGet() {
			params := strings.Split(path, "/")
			if len(params) < 2 || len(params) > 3 {
				ctx.Error("", fasthttp.StatusBadRequest)
				return
			}

			entity := params[0]
			idValue := params[1]

			if !(entity == "users" || entity == "locations" || entity == "visits") {
				ctx.Error("", fasthttp.StatusNotFound)
				return
			}

			id, err := strconv.Atoi(idValue)
			if err != nil {
				ctx.Error("", fasthttp.StatusNotFound)
				return
			}

			ctx.SetUserValue("id", int32(id))
			ctx.SetUserValue("entity", entity)

			if len(params) == 3 {
				if entity == "users" && params[2] == "visits" {
					handlers.UsersVisitsHandler(ctx)
					return
				}
				if entity == "locations" && params[2] == "avg" {
					handlers.LocationsAvgHandler(ctx)
					return
				}
				ctx.Error("", fasthttp.StatusBadRequest)
				return
			}

			handlers.EntitySelectHandler(ctx)
			return
		}

		ctx.Error("not found", fasthttp.StatusNotFound)
	}

	err := fasthttp.ListenAndServe(":"+os.Args[1], m)

	if err != nil {
		log.Fatal(err.Error())
	}
}
