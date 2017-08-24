package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"bytes"
	"fmt"
)

func UsersVisitsHandler(ctx *fasthttp.RequestCtx, id int64) {
	ctx.SetContentType("application/json;charset=utf-8")

	var fromDate, toDate, toDistance int
	var err error

	if id > storage.UserCount {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	user := storage.UserDb[id]
	if user == nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	args := ctx.QueryArgs()

	if args.Has("fromDate") {
		fromDate, err = args.GetUint("fromDate")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	if args.Has("toDate") {
		toDate, err = args.GetUint("toDate")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	if args.Has("toDistance") {
		toDistance, err = args.GetUint("toDistance")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	var country = (string)(args.Peek("country"))

	buffer := bufPool.Get().(*bytes.Buffer)
	buffer.Reset()
	buffer.WriteString(`{"visits": [`)

	atLeastOneFound := false
	for _, visit := range user.Visits {
		if fromDate > 0 && visit.Visited_at < int64(fromDate) {
			continue
		}
		if toDate > 0 && visit.Visited_at > int64(toDate) {
			continue
		}
		if len(country) > 0 && visit.Location_model.Country != country {
			continue
		}

		if toDistance > 0 && visit.Location_model.Distance >= int64(toDistance) {
			continue
		}

		if atLeastOneFound {
			buffer.WriteString(`,`)
		}

		buffer.WriteString(fmt.Sprintf("{\"mark\":%d,\"visited_at\":%d,\"place\":\"%s\"}", visit.Mark,visit.Visited_at,visit.Location_model.Place))
		atLeastOneFound = true
	}

	buffer.WriteString(`]}`)

	ctx.Write(buffer.Bytes())
	bufPool.Put(buffer)
}
