package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"encoding/json"
	"strconv"
)

func EntitySelectHandler(ctx *fasthttp.RequestCtx) {
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

	response, err := json.Marshal(entity)
	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(response)
}
