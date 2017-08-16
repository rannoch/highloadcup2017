package main

import (
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"fmt"
	"github.com/rannoch/highloadcup2017/storage"
)

func main() {
	storage.Init()

	LoadData("data")

	router := fasthttprouter.New()
	//router.POST("/:entity/new", entityNewHandler)
	router.GET("/:entity/:id", entitySelectHandler)
	router.POST("/:entity/:id", entityUpdateHandler)

	router.GET("/:entity/:id/visits", usersVisitsHandler)
	router.GET("/:entity/:id/avg", locationsAvgHandler)

	fasthttp.ListenAndServe(":8084", router.Handler)
}

func entitySelectHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s , %s!\n", ctx.UserValue("id"), ctx.UserValue("entity"))


}

func entityUpdateHandler(ctx *fasthttp.RequestCtx) {

}

func entityNewHandler(ctx *fasthttp.RequestCtx) {

}

func usersVisitsHandler(ctx *fasthttp.RequestCtx) {

}

func locationsAvgHandler(ctx *fasthttp.RequestCtx) {

}
