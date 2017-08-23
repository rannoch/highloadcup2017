package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"time"
	"fmt"
)

func LocationsAvgHandler(ctx *fasthttp.RequestCtx, id int32) {
	ctx.SetContentType("application/json;charset=utf-8")

	var fromDate, toDate, fromAge, toAge int
	var err error

	if id > storage.LocationCount {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	location := storage.LocationDb[id]
	if location == nil {
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
	if args.Has("fromAge") {
		fromAge, err = args.GetUint("fromAge")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}
	if args.Has("toAge") {
		toAge, err = args.GetUint("toAge")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	var gender = (string)(args.Peek("gender"))

	if gender != "" && !(gender == "m" || gender == "f") {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	var avg float32 = 0

	var marksSum int32 = 0
	var markCount int32 = 0

	for i := len(location.Visits) - 1; i >= 0; i -- {
		visit := location.Visits[i]
		if fromDate > 0 && visit.Visited_at <= int32(fromDate) {
			continue
		}
		if toDate > 0 && visit.Visited_at >= int32(toDate) {
			continue
		}
		if fromAge > 0 && visit.User_model.Birth_date >= int32(time.Now().AddDate(-fromAge, 0, 0).Unix()) {
			continue
		}
		if toAge > 0 && visit.User_model.Birth_date <= int32(time.Now().AddDate(-toAge, 0, 0).Unix()) {
			continue
		}
		if len(gender) > 0 && visit.User_model.Gender != gender {
			continue
		}

		marksSum += visit.Mark
		markCount++
	}

	if markCount > 0 {
		avg = float32(marksSum) / float32(markCount)
	}

	ctx.SetBodyString(fmt.Sprintf("{\"avg\" : %.5f}", avg))
}
