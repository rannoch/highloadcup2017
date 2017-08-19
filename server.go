package main

import (
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"github.com/rannoch/highloadcup2017/storage"
	"github.com/rannoch/highloadcup2017/handlers"
)


func main() {
	storage.Init()

	//LoadData("data")

	router := fasthttprouter.New()
	//router.POST("/:entity/new", entityNewHandler)
	router.GET("/:entity/:id", handlers.EntitySelectHandler)
	router.POST("/:entity/:id", handlers.EntityUpdateHandler)

	router.GET("/:entity/:id/visits", handlers.UsersVisitsHandler)
	router.GET("/:entity/:id/avg", handlers.LocationsAvgHandler)

	fasthttp.ListenAndServe(":8084", router.Handler)
}

