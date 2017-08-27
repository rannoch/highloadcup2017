package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/tcp_server/storage"
	"fmt"
	"github.com/rannoch/highloadcup2017/tcp_server/server"
	"bytes"
)

var (
	strM = []byte("m")
	strF = []byte("f")
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
	var gender []byte

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

		gender = ctx.QueryArgs.Peek("gender")

		if len(gender) > 0 && !(bytes.Equal(strM, gender) || bytes.Equal(strF, gender)) {
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
		if fromAge > 0 && visit.User_model.Birth_date >= storage.GenerateTime.AddDate(-fromAge, 0, 0).Unix() {
			continue
		}
		if toAge > 0 && visit.User_model.Birth_date <= storage.GenerateTime.AddDate(-toAge, 0, 0).Unix() {
			continue
		}
		if len(gender) > 0 && !bytes.Equal(visit.User_model.Gender, gender) {
			continue
		}

		marksSum += visit.Mark
		markCount++
	}

	if markCount > 0 {
		avg = float64(marksSum) / float64(markCount)
	}

	ctx.WriteString(fmt.Sprintf("{\"avg\" : %.5f}", avg))
}
