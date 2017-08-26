package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/tcp_server/storage"
	"time"
	"fmt"
	"github.com/rannoch/highloadcup2017/tcp_server/server"
	"bytes"
)

func LocationsAvgHandler(ctx *server.HlcupCtx, id int64) {
	var err error

	if id > storage.LocationCount {
		ctx.Error(fasthttp.StatusNotFound)
		return
	}

	location := storage.LocationDb[id]
	if location == nil {
		ctx.Error(fasthttp.StatusNotFound)
		return
	}

	var fromDate, toDate, fromAge, toAge int
	var gender string

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
		if ctx.QueryArgs.Has("fromAge") {
			fromAge, err = ctx.QueryArgs.GetUint("fromAge")

			if err != nil {
				ctx.Error(fasthttp.StatusBadRequest)
				return
			}
		}
		if ctx.QueryArgs.Has("toAge") {
			toAge, err = ctx.QueryArgs.GetUint("toAge")

			if err != nil {
				ctx.Error(fasthttp.StatusBadRequest)
				return
			}
		}

		gender = (string)(ctx.QueryArgs.Peek("gender"))

		if gender != "" && !(gender == "m" || gender == "f") {
			ctx.Error(fasthttp.StatusBadRequest)
			return
		}
	}

	var avg float64 = 0

	var marksSum int64 = 0
	var markCount int64 = 0

	for i := len(location.Visits) - 1; i >= 0; i -- {
		visit := location.Visits[i]
		if fromDate > 0 && visit.Visited_at <= int64(fromDate) {
			continue
		}
		if toDate > 0 && visit.Visited_at >= int64(toDate) {
			continue
		}
		if fromAge > 0 && visit.User_model.Birth_date >= time.Now().AddDate(-fromAge, 0, 0).Unix() {
			continue
		}
		if toAge > 0 && visit.User_model.Birth_date <= time.Now().AddDate(-toAge, 0, 0).Unix() {
			continue
		}
		if len(gender) > 0 && visit.User_model.Gender != gender {
			continue
		}

		marksSum += visit.Mark
		markCount++
	}

	if markCount > 0 {
		avg = float64(marksSum) / float64(markCount)
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("{\"avg\" : %.5f}", avg))

	ctx.Write(buffer.Bytes())
}
