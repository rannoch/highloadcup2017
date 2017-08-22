package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"github.com/rannoch/highloadcup2017/memory/models"
	"fmt"
	"bytes"
)

func UsersVisitsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id int32
	var fromDate, toDate, toDistance int
	var err error

	if ctx.QueryArgs().Has("fromDate") {
		fromDate, err = ctx.QueryArgs().GetUint("fromDate")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	if ctx.QueryArgs().Has("toDate") {
		toDate, err = ctx.QueryArgs().GetUint("toDate")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	if ctx.QueryArgs().Has("toDistance") {
		toDistance, err = ctx.QueryArgs().GetUint("toDistance")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	var country = (string)(ctx.QueryArgs().Peek("country"))

	id, _ = ctx.UserValue("id").(int32)

	if id > storage.UserCount {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	user := storage.UserDb[id]
	if user == nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	visits := []models.Visit{}

	for _, visit := range user.Visits {
		if fromDate > 0 && visit.Visited_at < int32(fromDate) {
			continue
		}
		if toDate > 0 && visit.Visited_at > int32(toDate) {
			continue
		}
		if len(country) > 0 && visit.Location_model.Country != country {
			continue
		}

		if toDistance > 0 && visit.Location_model.Distance >= int32(toDistance) {
			continue
		}

		visits = append(visits, *visit)
	}

	visitsResponse := ""
	for _, visit := range visits {
		visitsResponse += fmt.Sprintf("{\"mark\":%d,\"visited_at\":%d,\"place\":\"%s\"},", visit.Mark,visit.Visited_at,visit.Location_model.Place)
	}

	if len(visitsResponse) > 0 {
		visitsResponse = visitsResponse[:len(visitsResponse) - 1]
	}

	buffer := bufPool.Get().(*bytes.Buffer)
	buffer.Reset()
	buffer.WriteString("{\"visits\": [" + visitsResponse + "]}")

	ctx.Write(buffer.Bytes())
	bufPool.Put(buffer)
}
