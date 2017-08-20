package main

import (
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	//"github.com/rannoch/highloadcup2017/memory/models"
	"github.com/rannoch/highloadcup2017/memory/handlers"
	"os"
	"log"
	"github.com/rannoch/highloadcup2017/memory/storage"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("not enough args")
	}

	//storage.Init()

	storage.InitMemoryMap()
	LoadData(os.Args[2])


	router := fasthttprouter.New()
	//router.POST("/:entity/new", entityNewHandler)
	router.GET("/:entity/:id", handlers.EntitySelectHandler)
	router.POST("/:entity/:id", handlers.EntityUpdateHandler)

	router.GET("/:entity/:id/visits", handlers.UsersVisitsHandler)
	router.GET("/:entity/:id/avg", handlers.LocationsAvgHandler)

	err := fasthttp.ListenAndServe(":" + os.Args[1], router.Handler)

	if err != nil {
		log.Fatal(err.Error())
	}
}

