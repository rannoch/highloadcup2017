package handlers

import (
	"github.com/valyala/fasthttp"
	"strconv"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"encoding/json"
	"github.com/rannoch/highloadcup2017/memory/models"
)

func UsersVisitsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id int
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

	id, err = strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	u, exist := storage.Db["user"][int32(id)]
	if !exist {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	user := u.(*models.User)
	visits := user.Visits

	if fromDate > 0 {
		for i, visit := range visits {
			if visit.Visited_at < int32(fromDate) {
				visits = append(visits[:i], visits[i+1:]...)
			}
		}
	}

	if toDate > 0 {
		for i, visit := range visits {
			if visit.Visited_at > int32(toDate) {
				visits = append(visits[:i], visits[i+1:]...)
			}
		}
	}

	if len(country) > 0 {
		for i, visit := range visits {
			if visit.Location.Country != country {
				visits = append(visits[:i], visits[i+1:]...)
			}
		}
	}

	if toDistance > 0 {
		for i, visit := range visits {
			if visit.Location.Distance >= int32(toDistance) {
				visits = append(visits[:i], visits[i+1:]...)
			}
		}
	}

	visitsResponse := []interface{}{}

	for _, visit := range visits {
		v := map[string]interface{}{
			"mark":       visit.Mark,
			"visited_at": visit.Visited_at,
			"place":      visit.Location.Place,
		}

		visitsResponse = append(visitsResponse, v)
	}

	response, err := json.Marshal(map[string]interface{}{"visits": visitsResponse})
	if err != nil {
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(response)
}
