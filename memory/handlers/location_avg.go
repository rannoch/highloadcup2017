package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"github.com/rannoch/highloadcup2017/memory/models"
	"strconv"
	"encoding/json"
	"time"
)

func LocationsAvgHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id, fromDate, toDate, fromAge, toAge int
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

	id, err = strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	if gender != "" && !(gender == "m" || gender == "f") {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	l, exist := storage.Db["location"][int32(id)]
	if !exist {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	location := l.(*models.Location)
	visits := []models.Visit{}

	for _, visit := range location.Visits {
		visits = append(visits, *visit)
	}

	if fromDate > 0 {
		for i := len(visits) - 1; i >= 0; i -- {
			visit := visits[i]
			if visit.Visited_at < int32(fromDate) {
				visits = append(visits[:i], visits[i+1:]...)
			}
		}
	}

	if toDate > 0 {
		for i := len(visits) - 1; i >= 0; i -- {
			visit := visits[i]
			if visit.Visited_at > int32(toDate) {
				visits = append(visits[:i], visits[i+1:]...)
			}
		}
	}

	if fromAge > 0 {
		for i := len(visits) - 1; i >= 0; i -- {
			visit := visits[i]
			if visit.User_model.Birth_date >= int32(time.Now().AddDate(-fromAge, 0, 0).Unix()) {
				visits = append(visits[:i], visits[i+1:]...)
			}
		}
	}

	if toAge > 0 {
		for i := len(visits) - 1; i >= 0; i -- {
			visit := visits[i]
			if visit.User_model.Birth_date <= int32(time.Now().AddDate(-toAge, 0, 0).Unix()) {
				visits = append(visits[:i], visits[i+1:]...)
			}
		}
	}

	if len(gender) > 0 {
		for i := len(visits) - 1; i >= 0; i -- {
			visit := visits[i]
			if visit.User_model.Gender != gender {
				visits = append(visits[:i], visits[i+1:]...)
			}
		}
	}

	var marksSum int32 = 0
	for _, visit := range visits {
		marksSum += visit.Mark
	}

	if len(visits) > 0 {
		avg = float32(marksSum) / float32(len(visits))
	}

	response, err := json.Marshal(map[string]interface{}{"avg": models.FloatPrecision5(avg) })
	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(response)
}
