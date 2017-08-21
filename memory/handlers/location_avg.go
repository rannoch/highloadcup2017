package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"time"
	"fmt"
	"bytes"
)

func LocationsAvgHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var fromDate, toDate, fromAge, toAge int
	var id int32
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
	if ctx.QueryArgs().Has("fromAge") {
		fromAge, err = ctx.QueryArgs().GetUint("fromAge")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}
	if ctx.QueryArgs().Has("toAge") {
		toAge, err = ctx.QueryArgs().GetUint("toAge")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	var gender = (string)(ctx.QueryArgs().Peek("gender"))

	var avg float32 = 0

	id, _ = ctx.UserValue("id").(int32)

	if gender != "" && !(gender == "m" || gender == "f") {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	if id > storage.LocationCount {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	location := storage.LocationDb[id]
	if location == nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	var marksSum int32 = 0
	var markCount int32 = 0

	for i := len(location.Visits) - 1; i >= 0; i -- {
		visit := location.Visits[i]
		if fromDate > 0 && visit.Visited_at < int32(fromDate) {
			continue
		}
		if toDate > 0 && visit.Visited_at > int32(toDate) {
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

	//ctx.SetBody([]byte(fmt.Sprintf("{\"avg\" : %.5f}", avg)))

	buffer := bufPool.Get().(*bytes.Buffer)
	buffer.Reset()
	buffer.WriteString(fmt.Sprintf("{\"avg\" : %.5f}", avg))

	ctx.Write(buffer.Bytes())
	bufPool.Put(buffer)
}
