package handlers

import (
	//"github.com/valyala/fasthttp"
	//"github.com/rannoch/highloadcup2017/memory/storage"
	//"fmt"
	"github.com/rannoch/highloadcup2017/tcp_server/server"
	"github.com/valyala/fasthttp"
	"fmt"
	"github.com/rannoch/highloadcup2017/tcp_server/storage"
	"bytes"
	//"github.com/rannoch/highloadcup2017/tcp_server/logger"
)

func UsersVisitsHandler(ctx *server.HlcupCtx, id int64) {
	var fromDate, toDate, toDistance int
	var err error

	if id > storage.UserCount {
		ctx.Error(fasthttp.StatusNotFound)
		return
	}

	user := storage.UserDb[id]
	if user == nil {
		ctx.Error(fasthttp.StatusNotFound)
		return
	}

	if ctx.HasUrlParams {
		ctx.ParseParams()

		if ctx.QueryArgs.Has("fromDate") {
			fromDate, err = ctx.QueryArgs.GetUint("fromDate")

			if err != nil {
				ctx.Error(fasthttp.StatusBadRequest)
				return
			}
		}

		if ctx.QueryArgs.Has("toDate") {
			toDate, err = ctx.QueryArgs.GetUint("toDate")

			if err != nil {
				ctx.Error(fasthttp.StatusBadRequest)
				return
			}
		}

		if ctx.QueryArgs.Has("toDistance") {
			toDistance, err = ctx.QueryArgs.GetUint("toDistance")

			if err != nil {
				ctx.Error(fasthttp.StatusBadRequest)
				return
			}
		}
	}

	//logger.PrintLog(fmt.Sprintf("%d %d %d", fromDate, toDate, toDistance))

	var country = (string)(ctx.QueryArgs.Peek("country"))

	var buffer bytes.Buffer

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
}
